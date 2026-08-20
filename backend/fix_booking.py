import re

with open('internal/booking/service.go', 'r') as f:
    lines = f.readlines()

out = []
i = 0
while i < len(lines):
    line = lines[i]

    # Add design import
    if line.strip() == '"mehndi-booking-backend/internal/notification"':
        out.append(line)
        i += 1
        if i < len(lines) and lines[i].strip() == '"mehndi-booking-backend/internal/user"':
            out.append('\t"mehndi-booking-backend/internal/design"\n')
        continue

    # Add DesignGetter interface after AvailabilityRepository
    if line.strip() == '}' and i+1 < len(lines) and lines[i+1].strip() == '' and i+2 < len(lines) and lines[i+2].strip() == 'type Service struct {':
        out.append(line)
        out.append('\n')
        out.append('type DesignGetter interface {\n')
        out.append('\tGetDesignByID(ctx context.Context, id string) (design.Design, error)\n')
        out.append('}\n')
        i += 1
        continue

    # Update Service struct
    if line.strip() == 'type Service struct {':
        out.append('type Service struct {\n')
        out.append('\tbookingRepo      BookingRepository\n')
        out.append('\tavailabilityRepo AvailabilityRepository\n')
        out.append('\tnotificationSvc  notification.NotificationRepository\n')
        out.append('\tuserRepo         notification.UserRepository\n')
        out.append('\tdesignRepo       DesignGetter\n')
        out.append('}\n')
        while i < len(lines) and lines[i].strip() != '}':
            i += 1
        i += 1
        continue

    # Update NewService
    if line.strip() == 'func NewService(bookingRepo BookingRepository, availabilityRepo AvailabilityRepository, notificationSvc notification.NotificationRepository, userRepo notification.UserRepository) *Service {':
        out.append('func NewService(bookingRepo BookingRepository, availabilityRepo AvailabilityRepository, notificationSvc notification.NotificationRepository, userRepo notification.UserRepository, designRepo DesignGetter) *Service {\n')
        out.append('\treturn &Service{\n')
        out.append('\t\tbookingRepo:      bookingRepo,\n')
        out.append('\t\tavailabilityRepo: availabilityRepo,\n')
        out.append('\t\tnotificationSvc:  notificationSvc,\n')
        out.append('\t\tuserRepo:         userRepo,\n')
        out.append('\t\tdesignRepo:       designRepo,\n')
        out.append('\t}\n')
        out.append('}\n')
        while i < len(lines) and lines[i].strip() != '}':
            i += 1
        i += 1
        continue

    # Fix CreateBooking duration
    if line.strip() == 'endTime := startTime.Add(2 * time.Hour)':
        out.append('\tdesign, err := s.designRepo.GetDesignByID(ctx, req.DesignID)\n')
        out.append('\tif err != nil {\n')
        out.append('\t\treturn Booking{}, fmt.Errorf("failed to load design: %v", err)\n')
        out.append('\t}\n')
        out.append('\n')
        out.append('\tduration := time.Duration(design.DurationMinutes) * time.Minute\n')
        out.append('\tif duration <= 0 {\n')
        out.append('\t\tduration = 2 * time.Hour\n')
        out.append('\t}\n')
        out.append('\n')
        out.append('\tendTime := startTime.Add(duration)\n')
        i += 1
        continue

    # Fix CalculateAvailableSlots
    if line.strip() == 'var allSlots []Slot' and i+1 < len(lines) and lines[i+1].strip() == '' and i+2 < len(lines) and lines[i+2].strip().startswith('for _, av := range availabilities'):
        out.append('\tdesign, err := s.designRepo.GetDesignByID(ctx, designID)\n')
        out.append('\tif err != nil {\n')
        out.append('\t\treturn nil, fmt.Errorf("failed to load design: %v", err)\n')
        out.append('\t}\n')
        out.append('\tduration := time.Duration(design.DurationMinutes) * time.Minute\n')
        out.append('\tif duration <= 0 {\n')
        out.append('\t\tduration = 2 * time.Hour\n')
        out.append('\t}\n')
        out.append('\n')
        out.append('\tvar allSlots []Slot\n')
        out.append('\n')
        out.append('\tfor _, av := range availabilities {\n')
        out.append('\t\tslots := generateSlotsForAvailability(av, designID, bookingDateParsed, s.bookingRepo, duration)\n')
        out.append('\t\tallSlots = append(allSlots, slots...)\n')
        out.append('\t}\n')
        out.append('\n')
        out.append('\treturn allSlots, nil\n')
        out.append('}\n')
        out.append('\n')
        out.append('func generateSlotsForAvailability(av availability.Availability, designID string, bookingDate time.Time, bookingRepo BookingRepository, duration time.Duration) []Slot {\n')
        # skip old function until closing }
        while i < len(lines) and lines[i].strip() != '}':
            i += 1
        i += 1
        continue

    # Fix slot loop hardcoded 2*time.Hour
    if line.strip().startswith('for current.Add(2*time.Hour)'):
        out.append('for current.Add(duration).Before(endTime) || current.Add(duration).Equal(endTime) {\n')
        i += 1
        continue

    if line.strip().startswith('slotEnd := current.Add(2 * time.Hour)'):
        out.append('\t\tslotEnd := current.Add(duration)\n')
        i += 1
        continue

    out.append(line)
    i += 1

with open('internal/booking/service.go', 'w') as f:
    f.writelines(out)

print('Done')
