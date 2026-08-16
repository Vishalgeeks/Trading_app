package booking

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"mehndi-booking-backend/internal/availability"
	"mehndi-booking-backend/internal/notification"
	"mehndi-booking-backend/internal/user"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking Booking) (Booking, error)
	GetBookingByID(ctx context.Context, id string) (Booking, error)
	GetBookingByIDForUser(ctx context.Context, id string, userID string) (Booking, error)
	ListBookingsByUser(ctx context.Context, userID string, status *string, date *string) ([]Booking, error)
	ListAdminBookings(ctx context.Context, status *string, date *string, fromDate *string, toDate *string) ([]Booking, error)
	CheckBookingOverlap(ctx context.Context, designID string, bookingDate time.Time, startTime, endTime time.Time, excludeBookingID *string) (bool, error)
	UpdateBookingStatus(ctx context.Context, id string, status string) (Booking, error)
	CancelBooking(ctx context.Context, id string) (Booking, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type AvailabilityRepository interface {
	GetAvailabilityForDay(ctx context.Context, dayOfWeek int) ([]availability.Availability, error)
}

type Service struct {
	bookingRepo      BookingRepository
	availabilityRepo AvailabilityRepository
	notificationSvc  notification.NotificationRepository
	userRepo         notification.UserRepository
}

func NewService(bookingRepo BookingRepository, availabilityRepo AvailabilityRepository, notificationSvc notification.NotificationRepository, userRepo notification.UserRepository) *Service {
	return &Service{
		bookingRepo:      bookingRepo,
		availabilityRepo: availabilityRepo,
		notificationSvc:  notificationSvc,
		userRepo:         userRepo,
	}
}

func (s *Service) CreateBooking(ctx context.Context, userID string, req CreateBookingRequest) (Booking, error) {
	if req.DesignID == "" || req.BookingDate == "" || req.StartTime == "" {
		return Booking{}, fmt.Errorf("design_id, booking_date, and start_time are required")
	}

	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		return Booking{}, fmt.Errorf("invalid booking_date format, use YYYY-MM-DD")
	}

	if bookingDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return Booking{}, fmt.Errorf("booking_date cannot be in the past")
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return Booking{}, fmt.Errorf("invalid start_time format, use HH:MM")
	}

	dayOfWeek := int(bookingDate.Weekday())

	availabilities, err := s.availabilityRepo.GetAvailabilityForDay(ctx, dayOfWeek)
	if err != nil {
		return Booking{}, fmt.Errorf("failed to check availability: %v", err)
	}

	if len(availabilities) == 0 {
		return Booking{}, fmt.Errorf("no availability for selected date")
	}

	startStr := startTime.Format("15:04")
	var matchedAvail *availability.Availability
	for _, av := range availabilities {
		if startStr >= av.StartTime.Format("15:04") && startStr < av.EndTime.Format("15:04") {
			matchedAvail = &av
			break
		}
	}

	if matchedAvail == nil {
		return Booking{}, fmt.Errorf("selected time is outside admin availability")
	}

	endTime := startTime.Add(2 * time.Hour)

	if endTime.Format("15:04") > matchedAvail.EndTime.Format("15:04") {
		return Booking{}, fmt.Errorf("booking exceeds availability window")
	}

	tx, err := s.bookingRepo.(*Repository).pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return Booking{}, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	hasOverlap, err := s.checkOverlapInTx(ctx, tx, req.DesignID, bookingDate, startTime, endTime, nil)
	if err != nil {
		return Booking{}, fmt.Errorf("failed to check booking conflicts: %v", err)
	}

	if hasOverlap {
		return Booking{}, ErrBookingConflict
	}

	notes := ""
	if req.Notes != "" {
		notes = req.Notes
	}

	booking := Booking{
		UserID:      userID,
		DesignID:    req.DesignID,
		BookingDate: bookingDate,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "PENDING",
		Notes:       &notes,
	}

	created, err := s.createBookingInTx(ctx, tx, booking)
	if err != nil {
		if errors.Is(err, ErrBookingConflict) {
			return Booking{}, ErrBookingConflict
		}
		return Booking{}, fmt.Errorf("failed to create booking: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Booking{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Create notifications after successful commit
	if err := s.sendBookingCreatedNotifications(ctx, created); err != nil {
		slog.Error("failed to create booking notifications", "booking_id", created.ID, "error", err)
	}

	return created, nil
}

func (s *Service) checkOverlapInTx(ctx context.Context, tx pgx.Tx, designID string, bookingDate time.Time, startTime, endTime time.Time, excludeBookingID *string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM bookings
			WHERE design_id = $1
			  AND booking_date = $2::date
			  AND status IN ('PENDING', 'CONFIRMED')
			  AND start_time < $4
			  AND end_time > $3
		`
	args := []interface{}{designID, bookingDate.Format("2006-01-02"), startTime, endTime}
	argCount := 4

	if excludeBookingID != nil && *excludeBookingID != "" {
		query += fmt.Sprintf(" AND id != $%d", argCount+1)
		args = append(args, *excludeBookingID)
	}

	query += ")"

	var exists bool
	row := tx.QueryRow(ctx, query, args...)
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return exists, nil
}

func (s *Service) createBookingInTx(ctx context.Context, tx pgx.Tx, booking Booking) (Booking, error) {
	query := `
		INSERT INTO bookings (user_id, design_id, booking_date, start_time, end_time, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, design_id, booking_date, start_time, end_time, status, notes, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		booking.UserID,
		booking.DesignID,
		booking.BookingDate,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		booking.Notes,
	)

	created, err := scanBooking(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return Booking{}, ErrBookingConflict
			}
		}
		return Booking{}, fmt.Errorf("%w: %v", ErrDatabaseFailure, err)
	}

	return created, nil
}

func (s *Service) GetBooking(ctx context.Context, id string, userID string, isAdmin bool) (Booking, error) {
	if isAdmin {
		booking, err := s.bookingRepo.GetBookingByID(ctx, id)
		if err != nil {
			return Booking{}, err
		}
		return booking, nil
	}

	booking, err := s.bookingRepo.GetBookingByIDForUser(ctx, id, userID)
	if err != nil {
		return Booking{}, err
	}
	return booking, nil
}

func (s *Service) ListUserBookings(ctx context.Context, userID string, status *string, date *string) ([]Booking, error) {
	return s.bookingRepo.ListBookingsByUser(ctx, userID, status, date)
}

func (s *Service) ListAdminBookings(ctx context.Context, status *string, date *string, fromDate *string, toDate *string) ([]Booking, error) {
	return s.bookingRepo.ListAdminBookings(ctx, status, date, fromDate, toDate)
}

func (s *Service) CancelBooking(ctx context.Context, id string, userID string, isAdmin bool) (Booking, error) {
	if isAdmin {
		booking, err := s.bookingRepo.CancelBooking(ctx, id)
		if err != nil {
			return Booking{}, err
		}
		if err := s.notifyCancellation(ctx, booking, false); err != nil {
			slog.Error("failed to create cancellation notification", "booking_id", booking.ID, "error", err)
		}
		return booking, nil
	}

	booking, err := s.bookingRepo.GetBookingByIDForUser(ctx, id, userID)
	if err != nil {
		return Booking{}, err
	}

	if booking.Status == "CANCELLED" || booking.Status == "COMPLETED" {
		return Booking{}, fmt.Errorf("cannot cancel a booking with status: %s", booking.Status)
	}

	updated, err := s.bookingRepo.CancelBooking(ctx, id)
	if err != nil {
		return Booking{}, err
	}

	if err := s.notifyCancellation(ctx, updated, true); err != nil {
		slog.Error("failed to create cancellation notification", "booking_id", updated.ID, "error", err)
	}

	return updated, nil
}

func (s *Service) notifyCancellation(ctx context.Context, b Booking, cancelledByClient bool) error {
	if s.notificationSvc == nil || s.userRepo == nil {
		return nil
	}

	if cancelledByClient {
		// Notify admin(s)
		admins, err := s.userRepo.ListUsersByRole(ctx, user.RoleAdmin)
		if err != nil {
			return fmt.Errorf("failed to get admin users: %w", err)
		}
		for _, admin := range admins {
			adminNotification := notification.Notification{
				UserID:    admin.ID,
				Type:      notification.TypeClientCancelledBooking,
				Title:     "Booking Cancelled by Client",
				Message:   fmt.Sprintf("A client has cancelled the booking for %s on %s at %s.", b.DesignName, b.BookingDate.Format("2006-01-02"), b.StartTime.Format("15:04")),
				BookingID: &b.ID,
				IsRead:    false,
			}
			_, err := s.notificationSvc.CreateNotification(ctx, adminNotification)
			if err != nil {
				return fmt.Errorf("failed to create admin cancellation notification: %w", err)
			}
		}
	} else {
		// Notify client
		clientNotification := notification.Notification{
			UserID:    b.UserID,
			Type:      notification.TypeBookingCancelled,
			Title:     "Booking Cancelled",
			Message:   fmt.Sprintf("Your booking for %s on %s at %s has been cancelled.", b.DesignName, b.BookingDate.Format("2006-01-02"), b.StartTime.Format("15:04")),
			BookingID: &b.ID,
			IsRead:    false,
		}
		_, err := s.notificationSvc.CreateNotification(ctx, clientNotification)
		if err != nil {
			return fmt.Errorf("failed to create client cancellation notification: %w", err)
		}
	}
	return nil
}

func (s *Service) UpdateBookingStatus(ctx context.Context, id string, newStatus string) (Booking, error) {
	validTransitions := map[string]bool{
		"PENDING":   true,
		"CONFIRMED": true,
		"CANCELLED": true,
		"COMPLETED": true,
	}

	if !validTransitions[newStatus] {
		return Booking{}, ErrInvalidStatus
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		return Booking{}, err
	}

	oldStatus := booking.Status

	switch booking.Status {
	case "PENDING":
		if newStatus != "CONFIRMED" && newStatus != "CANCELLED" {
			return Booking{}, fmt.Errorf("cannot transition from PENDING to %s", newStatus)
		}
	case "CONFIRMED":
		if newStatus != "CANCELLED" && newStatus != "COMPLETED" {
			return Booking{}, fmt.Errorf("cannot transition from CONFIRMED to %s", newStatus)
		}
	case "CANCELLED":
		return Booking{}, fmt.Errorf("cannot update a cancelled booking")
	case "COMPLETED":
		return Booking{}, fmt.Errorf("cannot update a completed booking")
	default:
		return Booking{}, ErrInvalidStatus
	}

	updated, err := s.bookingRepo.UpdateBookingStatus(ctx, id, newStatus)
	if err != nil {
		return Booking{}, err
	}

	// Send notifications based on transition
	if err := s.sendStatusChangeNotification(ctx, updated, oldStatus, newStatus); err != nil {
		slog.Error("failed to create status notification", "booking_id", updated.ID, "error", err)
	}

	return updated, nil
}

func (s *Service) sendStatusChangeNotification(ctx context.Context, b Booking, oldStatus, newStatus string) error {
	switch {
	case oldStatus == "PENDING" && newStatus == "CONFIRMED":
		return s.notifyConfirmation(ctx, b)
	case oldStatus == "CONFIRMED" && newStatus == "COMPLETED":
		return s.notifyCompletion(ctx, b)
	}
	return nil
}

func (s *Service) notifyConfirmation(ctx context.Context, b Booking) error {
	if s.notificationSvc == nil {
		return nil
	}

	// Check for duplicate
	existing, err := s.notificationSvc.ListNotificationsByUser(ctx, b.UserID, 10, 0)
	if err == nil {
		for _, n := range existing {
			if n.BookingID != nil && *n.BookingID == b.ID && n.Type == notification.TypeBookingConfirmed {
				return nil
			}
		}
	}

	clientNotification := notification.Notification{
		UserID:    b.UserID,
		Type:      notification.TypeBookingConfirmed,
		Title:     "Booking Confirmed",
		Message:   fmt.Sprintf("Your booking for %s on %s at %s has been confirmed.", b.DesignName, b.BookingDate.Format("2006-01-02"), b.StartTime.Format("15:04")),
		BookingID: &b.ID,
		IsRead:    false,
	}
	_, err = s.notificationSvc.CreateNotification(ctx, clientNotification)
	if err != nil {
		return fmt.Errorf("failed to create confirmation notification: %w", err)
	}
	return nil
}

func (s *Service) notifyCompletion(ctx context.Context, b Booking) error {
	if s.notificationSvc == nil {
		return nil
	}

	// Check for duplicate
	existing, err := s.notificationSvc.ListNotificationsByUser(ctx, b.UserID, 10, 0)
	if err == nil {
		for _, n := range existing {
			if n.BookingID != nil && *n.BookingID == b.ID && n.Type == notification.TypeBookingCompleted {
				return nil
			}
		}
	}

	clientNotification := notification.Notification{
		UserID:    b.UserID,
		Type:      notification.TypeBookingCompleted,
		Title:     "Booking Completed",
		Message:   fmt.Sprintf("Your booking for %s on %s at %s has been marked as completed.", b.DesignName, b.BookingDate.Format("2006-01-02"), b.StartTime.Format("15:04")),
		BookingID: &b.ID,
		IsRead:    false,
	}
	_, err = s.notificationSvc.CreateNotification(ctx, clientNotification)
	if err != nil {
		return fmt.Errorf("failed to create completion notification: %w", err)
	}
	return nil
}

func (s *Service) CalculateAvailableSlots(ctx context.Context, designID string, bookingDate string) ([]Slot, error) {
	bookingDateParsed, err := time.Parse("2006-01-02", bookingDate)
	if err != nil {
		return nil, fmt.Errorf("invalid booking_date format, use YYYY-MM-DD")
	}

	dayOfWeek := int(bookingDateParsed.Weekday())

	availabilities, err := s.availabilityRepo.GetAvailabilityForDay(ctx, dayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get availability: %v", err)
	}

	if len(availabilities) == 0 {
		return nil, fmt.Errorf("no availability for selected date")
	}

	var allSlots []Slot

	for _, av := range availabilities {
		slots := generateSlotsForAvailability(av, designID, bookingDateParsed, s.bookingRepo)
		allSlots = append(allSlots, slots...)
	}

	return allSlots, nil
}

func generateSlotsForAvailability(av availability.Availability, designID string, bookingDate time.Time, bookingRepo BookingRepository) []Slot {
	var slots []Slot

	startTime, _ := time.Parse("15:04", av.StartTime.Format("15:04"))
	endTime, _ := time.Parse("15:04", av.EndTime.Format("15:04"))

	interval := 30 * time.Minute
	current := startTime

	for current.Add(2*time.Hour).Before(endTime) || current.Add(2*time.Hour).Equal(endTime) {
		slotEnd := current.Add(2 * time.Hour)

		hasOverlap, _ := bookingRepo.CheckBookingOverlap(
			context.Background(),
			designID,
			bookingDate,
			current,
			slotEnd,
			nil,
		)

		if !hasOverlap {
			slots = append(slots, Slot{
				StartTime: current.Format("15:04"),
				EndTime:   slotEnd.Format("15:04"),
			})
		}

		current = current.Add(interval)
	}

	return slots
}

func (s *Service) sendBookingCreatedNotifications(ctx context.Context, b Booking) error {
	if s.notificationSvc == nil || s.userRepo == nil {
		return nil
	}

	// Notify client
	clientNotification := notification.Notification{
		UserID:    b.UserID,
		Type:      notification.TypeBookingCreated,
		Title:     "Booking Request Submitted",
		Message:   fmt.Sprintf("Your booking request for %s on %s at %s has been submitted.", b.DesignName, b.BookingDate.Format("2006-01-02"), b.StartTime.Format("15:04")),
		BookingID: &b.ID,
		IsRead:    false,
	}
	_, err := s.notificationSvc.CreateNotification(ctx, clientNotification)
	if err != nil {
		return fmt.Errorf("failed to create client notification: %w", err)
	}

	// Notify admin(s)
	admins, err := s.userRepo.ListUsersByRole(ctx, user.RoleAdmin)
	if err != nil {
		return fmt.Errorf("failed to get admin users: %w", err)
	}

	for _, admin := range admins {
		adminNotification := notification.Notification{
			UserID:    admin.ID,
			Type:      notification.TypeNewBooking,
			Title:     "New Booking Received",
			Message:   fmt.Sprintf("A new booking has been received for %s on %s at %s.", b.DesignName, b.BookingDate.Format("2006-01-02"), b.StartTime.Format("15:04")),
			BookingID: &b.ID,
			IsRead:    false,
		}
		_, err := s.notificationSvc.CreateNotification(ctx, adminNotification)
		if err != nil {
			return fmt.Errorf("failed to create admin notification: %w", err)
		}
	}

	return nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(i int) *int {
	return &i
}
