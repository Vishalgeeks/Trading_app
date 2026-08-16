package design

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"mehndi-booking-backend/internal/category"
)

type CategoryGetter interface {
	GetCategory(ctx context.Context, id string) (category.Category, error)
}

type DesignRepository interface {
	CreateDesign(ctx context.Context, des Design) (Design, error)
	GetDesignByID(ctx context.Context, id string) (Design, error)
	GetDesignBySlug(ctx context.Context, slug string) (Design, error)
	ListDesigns(ctx context.Context, categoryID *string, searchQuery *string, page, limit int) ([]Design, error)
	ListActiveDesigns(ctx context.Context) ([]Design, error)
	ListDesignsByCategory(ctx context.Context, categoryID string) ([]Design, error)
	SearchDesigns(ctx context.Context, query string) ([]Design, error)
	UpdateDesign(ctx context.Context, id string, req UpdateDesignRequest) (Design, error)
	DeactivateDesign(ctx context.Context, id string) error
}

type Service struct {
	repo         DesignRepository
	categoryRepo CategoryGetter
}

func NewService(repo DesignRepository, categoryRepo CategoryGetter) *Service {
	return &Service{repo: repo, categoryRepo: categoryRepo}
}

func (s *Service) CreateDesign(ctx context.Context, req CreateDesignRequest) (Design, error) {
	if strings.TrimSpace(req.Name) == "" {
		return Design{}, ErrInvalidName
	}
	if len(strings.TrimSpace(req.Name)) > 150 {
		return Design{}, fmt.Errorf("name exceeds maximum length of 150 characters")
	}
	if strings.TrimSpace(req.Slug) == "" {
		return Design{}, ErrInvalidSlug
	}
	if !ValidateSlug(req.Slug) {
		return Design{}, fmt.Errorf("invalid slug format")
	}
	if strings.TrimSpace(req.ImageURL) == "" {
		return Design{}, ErrInvalidImageURL
	}
	if req.Price == "" {
		return Design{}, ErrInvalidPrice
	}

	price := strings.TrimSpace(req.Price)
	if price[0] == '-' {
		return Design{}, ErrInvalidPrice
	}

	if req.DurationMinutes <= 0 {
		return Design{}, ErrInvalidDuration
	}

	cat, err := s.categoryRepo.GetCategory(ctx, req.CategoryID)
	if err != nil {
		if errors.Is(err, category.ErrCategoryNotFound) {
			return Design{}, ErrCategoryNotFound
		}
		return Design{}, err
	}

	if !cat.IsActive {
		return Design{}, ErrCategoryInactive
	}

	slug := strings.ToLower(strings.TrimSpace(req.Slug))

	des := Design{
		CategoryID:      req.CategoryID,
		Name:            strings.TrimSpace(req.Name),
		Slug:            slug,
		Description:     req.Description,
		ImageURL:        strings.TrimSpace(req.ImageURL),
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		IsActive:        true,
	}

	created, err := s.repo.CreateDesign(ctx, des)
	if err != nil {
		return Design{}, err
	}

	return created, nil
}

func (s *Service) GetDesign(ctx context.Context, id string) (Design, error) {
	if id == "" {
		return Design{}, ErrDesignNotFound
	}
	return s.repo.GetDesignByID(ctx, id)
}

func (s *Service) ListDesigns(ctx context.Context, categoryID *string, searchQuery *string, page, limit int) ([]Design, error) {
	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	return s.repo.ListDesigns(ctx, categoryID, searchQuery, page, limit)
}

func (s *Service) ListActiveDesigns(ctx context.Context) ([]Design, error) {
	return s.repo.ListActiveDesigns(ctx)
}

func (s *Service) ListDesignsByCategory(ctx context.Context, categoryID string) ([]Design, error) {
	return s.repo.ListDesignsByCategory(ctx, categoryID)
}

func (s *Service) SearchDesigns(ctx context.Context, query string) ([]Design, error) {
	return s.repo.SearchDesigns(ctx, query)
}

func (s *Service) UpdateDesign(ctx context.Context, id string, req UpdateDesignRequest) (Design, error) {
	if id == "" {
		return Design{}, ErrDesignNotFound
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return Design{}, ErrInvalidName
		}
		if utf8.RuneCountInString(name) > 150 {
			return Design{}, fmt.Errorf("name exceeds maximum length of 150 characters")
		}
		req.Name = &name
	}

	if req.Slug != nil {
		slug := strings.ToLower(strings.TrimSpace(*req.Slug))
		if slug == "" {
			return Design{}, ErrInvalidSlug
		}
		if !ValidateSlug(slug) {
			return Design{}, fmt.Errorf("invalid slug format")
		}
		req.Slug = &slug
	}

	if req.ImageURL != nil {
		imageURL := strings.TrimSpace(*req.ImageURL)
		if imageURL == "" {
			return Design{}, ErrInvalidImageURL
		}
		req.ImageURL = &imageURL
	}

	if req.CategoryID != nil {
		cat, err := s.categoryRepo.GetCategory(ctx, *req.CategoryID)
		if err != nil {
			if errors.Is(err, category.ErrCategoryNotFound) {
				return Design{}, ErrCategoryNotFound
			}
			return Design{}, err
		}
		if !cat.IsActive {
			return Design{}, ErrCategoryInactive
		}
	}

	des, err := s.repo.UpdateDesign(ctx, id, req)
	if err != nil {
		return Design{}, err
	}

	return des, nil
}

func (s *Service) DeactivateDesign(ctx context.Context, id string) error {
	if id == "" {
		return ErrDesignNotFound
	}
	return s.repo.DeactivateDesign(ctx, id)
}
