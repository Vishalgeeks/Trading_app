package favorite

import (
	"context"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateFavorite(ctx context.Context, userID, designID string) (Favorite, error) {
	if userID == "" || designID == "" {
		return Favorite{}, errors.New("user_id and design_id are required")
	}

	return s.repo.CreateFavorite(ctx, Favorite{
		UserID:   userID,
		DesignID: designID,
	})
}

func (s *Service) RemoveFavorite(ctx context.Context, userID, designID string) error {
	if userID == "" || designID == "" {
		return errors.New("user_id and design_id are required")
	}

	return s.repo.DeleteFavorite(ctx, userID, designID)
}

func (s *Service) GetUserFavorites(ctx context.Context, userID string) ([]Favorite, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	return s.repo.ListFavoritesByUser(ctx, userID)
}

func (s *Service) IsFavorite(ctx context.Context, userID, designID string) (bool, error) {
	if userID == "" || designID == "" {
		return false, nil
	}

	_, err := s.repo.GetFavorite(ctx, userID, designID)
	if err != nil {
		return false, nil
	}

	return true, nil
}
