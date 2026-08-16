package availability

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAvailability(ctx context.Context, av Availability) (Availability, error) {
	if av.DayOfWeek < 0 || av.DayOfWeek > 6 {
		return Availability{}, fmt.Errorf("day_of_week must be between 0 and 6")
	}

	startTime, err := time.Parse("15:04", av.StartTime.Format("15:04"))
	if err != nil {
		return Availability{}, fmt.Errorf("invalid start_time format")
	}
	endTime, err := time.Parse("15:04", av.EndTime.Format("15:04"))
	if err != nil {
		return Availability{}, fmt.Errorf("invalid end_time format")
	}

	if !endTime.After(startTime) {
		return Availability{}, fmt.Errorf("end_time must be after start_time")
	}

	av.StartTime = startTime
	av.EndTime = endTime

	return s.repo.CreateAvailability(ctx, av)
}

func (s *Service) GetAvailability(ctx context.Context, id string) (Availability, error) {
	return s.repo.GetAvailabilityByID(ctx, id)
}

func (s *Service) ListAvailability(ctx context.Context, activeOnly bool) ([]Availability, error) {
	return s.repo.ListAvailability(ctx, activeOnly)
}

func (s *Service) GetAvailabilityForDay(ctx context.Context, dayOfWeek int) ([]Availability, error) {
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return nil, fmt.Errorf("day_of_week must be between 0 and 6")
	}
	return s.repo.GetAvailabilityForDay(ctx, dayOfWeek)
}

func (s *Service) UpdateAvailability(ctx context.Context, id string, req UpdateAvailabilityRequest) (Availability, error) {
	if req.StartTime != nil {
		startTime, err := time.Parse("15:04", *req.StartTime)
		if err != nil {
			return Availability{}, fmt.Errorf("invalid start_time format")
		}
		startStr := startTime.Format("15:04")
		req.StartTime = &startStr
	}

	if req.EndTime != nil {
		endTime, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			return Availability{}, fmt.Errorf("invalid end_time format")
		}
		endStr := endTime.Format("15:04")
		req.EndTime = &endStr
	}

	av, err := s.repo.UpdateAvailability(ctx, id, req)
	if err != nil {
		if errors.Is(err, ErrAvailabilityNotFound) {
			return Availability{}, ErrAvailabilityNotFound
		}
		return Availability{}, err
	}

	return av, nil
}

func (s *Service) DeactivateAvailability(ctx context.Context, id string) error {
	return s.repo.DeactivateAvailability(ctx, id)
}

func (s *Service) DeleteAvailability(ctx context.Context, id string) error {
	return s.repo.DeleteAvailability(ctx, id)
}
