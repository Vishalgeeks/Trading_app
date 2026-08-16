package booking

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"mehndi-booking-backend/internal/availability"
)

var (
	testPool *pgxpool.Pool
	ctx      = context.Background()
)

func getPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testPool == nil {
		pool, err := pgxpool.New(ctx, "postgres://postgres:masum7003@localhost:5432/mehndi_booking?sslmode=disable")
		require.NoError(t, err)
		testPool = pool
	}
	return testPool
}

func setupBookingRepo(t *testing.T) (*Repository, *availability.Repository) {
	t.Helper()
	pool := getPool(t)
	return NewRepository(pool), availability.NewRepository(pool)
}

func createTestAvailability(t *testing.T, avRepo *availability.Repository, bookingDate string) availability.Availability {
	t.Helper()
	date, _ := time.Parse("2006-01-02", bookingDate)
	dayOfWeek := int(date.Weekday())
	startTime, _ := time.Parse("15:04", "09:00")
	endTime, _ := time.Parse("15:04", "17:00")
	av, err := avRepo.CreateAvailability(ctx, availability.Availability{
		DayOfWeek: dayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  true,
	})
	require.NoError(t, err)
	return av
}

func TestMain(m *testing.M) {
	pool, err := pgxpool.New(ctx, "postgres://postgres:masum7003@localhost:5432/mehndi_booking?sslmode=disable")
	if err != nil {
		panic(err)
	}
	testPool = pool
	pool.Exec(ctx, "DELETE FROM bookings WHERE id != '00000000-0000-0000-0000-000000000000'")
	pool.Exec(ctx, "DELETE FROM admin_availability WHERE id != '00000000-0000-0000-0000-000000000000'")
	pool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash, role, is_active) VALUES ('00000000-0000-0000-0000-000000000001', 'Test User 1', 'test1@example.com', 'hash', 'CLIENT', true) ON CONFLICT (id) DO NOTHING")
	pool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash, role, is_active) VALUES ('00000000-0000-0000-0000-000000000002', 'Test User 2', 'test2@example.com', 'hash', 'CLIENT', true) ON CONFLICT (id) DO NOTHING")
	pool.Exec(ctx, "INSERT INTO users (id, name, email, password_hash, role, is_active) VALUES ('00000000-0000-0000-0000-000000000003', 'Test User 3', 'test3@example.com', 'hash', 'CLIENT', true) ON CONFLICT (id) DO NOTHING")
	pool.Exec(ctx, "INSERT INTO categories (id, name, slug, is_active) VALUES ('00000000-0000-0000-0000-000000000004', 'Test Category', 'test-category', true) ON CONFLICT (id) DO NOTHING")
	pool.Exec(ctx, "INSERT INTO designs (id, category_id, name, slug, image_url, price, duration_minutes, is_active) VALUES ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000004', 'Test Design', 'test-design', 'http://example.com/test.jpg', '100.00', 120, true) ON CONFLICT (id) DO NOTHING")

	m.Run()

	pool.Exec(ctx, "DELETE FROM bookings WHERE id != '00000000-0000-0000-0000-000000000000'")
	pool.Exec(ctx, "DELETE FROM admin_availability WHERE id != '00000000-0000-0000-0000-000000000000'")
	pool.Exec(ctx, "DELETE FROM users WHERE id IN ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003')")
	pool.Exec(ctx, "DELETE FROM designs WHERE id = '00000000-0000-0000-0000-000000000002'")
	pool.Exec(ctx, "DELETE FROM categories WHERE id = '00000000-0000-0000-0000-000000000004'")
}

func TestRepository_CreateBooking(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-25")

	bookingDate, _ := time.Parse("2006-01-02", "2026-08-25")
	startTime, _ := time.Parse("15:04", "10:00")
	endTime, _ := time.Parse("15:04", "12:00")

	booking, err := bookingRepo.CreateBooking(ctx, Booking{
		UserID:      "00000000-0000-0000-0000-000000000001",
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: bookingDate,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "PENDING",
	})
	require.NoError(t, err)
	require.Equal(t, "00000000-0000-0000-0000-000000000001", booking.UserID)
	require.Equal(t, "00000000-0000-0000-0000-000000000002", booking.DesignID)
}

func TestRepository_GetBookingByID(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-25")

	bookingDate, _ := time.Parse("2006-01-02", "2026-08-25")
	startTime, _ := time.Parse("15:04", "10:00")
	endTime, _ := time.Parse("15:04", "12:00")

	created, err := bookingRepo.CreateBooking(ctx, Booking{
		UserID:      "00000000-0000-0000-0000-000000000001",
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: bookingDate,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "PENDING",
	})
	require.NoError(t, err)

	found, err := bookingRepo.GetBookingByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, found.ID)
}

func TestRepository_CheckBookingOverlap(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-25")

	bookingDate, _ := time.Parse("2006-01-02", "2026-08-25")
	startTime, _ := time.Parse("15:04", "10:00")
	endTime, _ := time.Parse("15:04", "12:00")

	_, err := bookingRepo.CreateBooking(ctx, Booking{
		UserID:      "00000000-0000-0000-0000-000000000001",
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: bookingDate,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "PENDING",
	})
	require.NoError(t, err)

	overlapStart, _ := time.Parse("15:04", "11:00")
	overlapEnd, _ := time.Parse("15:04", "13:00")
	hasOverlap, err := bookingRepo.CheckBookingOverlap(ctx, "00000000-0000-0000-0000-000000000002", bookingDate, overlapStart, overlapEnd, nil)
	require.NoError(t, err)
	require.True(t, hasOverlap)

	noOverlapStart, _ := time.Parse("15:04", "12:00")
	noOverlapEnd, _ := time.Parse("15:04", "14:00")
	hasOverlap2, err := bookingRepo.CheckBookingOverlap(ctx, "00000000-0000-0000-0000-000000000002", bookingDate, noOverlapStart, noOverlapEnd, nil)
	require.NoError(t, err)
	require.False(t, hasOverlap2)
}

func TestRepository_UpdateBookingStatus(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-25")

	bookingDate, _ := time.Parse("2006-01-02", "2026-08-25")
	startTime, _ := time.Parse("15:04", "10:00")
	endTime, _ := time.Parse("15:04", "12:00")

	created, err := bookingRepo.CreateBooking(ctx, Booking{
		UserID:      "00000000-0000-0000-0000-000000000001",
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: bookingDate,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "PENDING",
	})
	require.NoError(t, err)

	updated, err := bookingRepo.UpdateBookingStatus(ctx, created.ID, "CONFIRMED")
	require.NoError(t, err)
	require.Equal(t, "CONFIRMED", updated.Status)
}

func TestRepository_CancelBooking(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-25")

	bookingDate, _ := time.Parse("2006-01-02", "2026-08-25")
	startTime, _ := time.Parse("15:04", "10:00")
	endTime, _ := time.Parse("15:04", "12:00")

	created, err := bookingRepo.CreateBooking(ctx, Booking{
		UserID:      "00000000-0000-0000-0000-000000000001",
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: bookingDate,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "PENDING",
	})
	require.NoError(t, err)

	cancelled, err := bookingRepo.CancelBooking(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "CANCELLED", cancelled.Status)
}

func TestService_CreateBooking_Valid(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-26")
	svc := NewService(bookingRepo, avRepo, nil, nil)

	booking, err := svc.CreateBooking(ctx, "00000000-0000-0000-0000-000000000001", CreateBookingRequest{
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: "2026-08-26",
		StartTime:   "10:00",
	})
	require.NoError(t, err)
	require.Equal(t, "PENDING", booking.Status)
}

func TestService_CreateBooking_Overlap(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-27")
	svc := NewService(bookingRepo, avRepo, nil, nil)

	_, err := svc.CreateBooking(ctx, "00000000-0000-0000-0000-000000000001", CreateBookingRequest{
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: "2026-08-27",
		StartTime:   "10:00",
	})
	require.NoError(t, err)

	_, err = svc.CreateBooking(ctx, "00000000-0000-0000-0000-000000000003", CreateBookingRequest{
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: "2026-08-27",
		StartTime:   "11:00",
	})
	require.Error(t, err)
}

func TestService_CancelBooking_Valid(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-28")
	svc := NewService(bookingRepo, avRepo, nil, nil)

	booking, err := svc.CreateBooking(ctx, "00000000-0000-0000-0000-000000000001", CreateBookingRequest{
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: "2026-08-28",
		StartTime:   "10:00",
	})
	require.NoError(t, err)

	cancelled, err := svc.CancelBooking(ctx, booking.ID, "00000000-0000-0000-0000-000000000001", false)
	require.NoError(t, err)
	require.Equal(t, "CANCELLED", cancelled.Status)
}

func TestService_UpdateBookingStatus_Valid(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-29")
	svc := NewService(bookingRepo, avRepo, nil, nil)

	booking, err := svc.CreateBooking(ctx, "00000000-0000-0000-0000-000000000001", CreateBookingRequest{
		DesignID:    "00000000-0000-0000-0000-000000000002",
		BookingDate: "2026-08-29",
		StartTime:   "10:00",
	})
	require.NoError(t, err)

	updated, err := svc.UpdateBookingStatus(ctx, booking.ID, "CONFIRMED")
	require.NoError(t, err)
	require.Equal(t, "CONFIRMED", updated.Status)
}

func TestConcurrentBooking(t *testing.T) {
	bookingRepo, avRepo := setupBookingRepo(t)
	_ = createTestAvailability(t, avRepo, "2026-08-30")
	svc := NewService(bookingRepo, avRepo, nil, nil)

	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func(userID string) {
			_, err := svc.CreateBooking(ctx, userID, CreateBookingRequest{
				DesignID:    "00000000-0000-0000-0000-000000000002",
				BookingDate: "2026-08-30",
				StartTime:   "10:00",
			})
			results <- err
		}(fmt.Sprintf("00000000-0000-0000-0000-00000000000%d", i+1))
	}

	err1 := <-results
	err2 := <-results

	t.Logf("err1: %v, err2: %v", err1, err2)
	require.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil), "exactly one booking should succeed")
}
