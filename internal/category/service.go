package category

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, cat Category) (Category, error)
	GetCategoryByID(ctx context.Context, id string) (Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (Category, error)
	ListCategories(ctx context.Context, activeOnly bool) ([]Category, error)
	UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (Category, error)
	DeactivateCategory(ctx context.Context, id string) error
	CountActiveDesignsInCategory(ctx context.Context, categoryID string) (int, error)
}

type Service struct {
	repo CategoryRepository
}

func NewService(repo CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, req CreateCategoryRequest) (Category, error) {
	if strings.TrimSpace(req.Name) == "" {
		return Category{}, ErrInvalidName
	}
	if len(strings.TrimSpace(req.Name)) > 100 {
		return Category{}, fmt.Errorf("name exceeds maximum length of 100 characters")
	}
	if strings.TrimSpace(req.Slug) == "" {
		return Category{}, ErrInvalidSlug
	}
	if !ValidateSlug(req.Slug) {
		return Category{}, fmt.Errorf("invalid slug format")
	}

	slug := strings.ToLower(strings.TrimSpace(req.Slug))

	cat := Category{
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Description: req.Description,
		IsActive:    true,
	}

	created, err := s.repo.CreateCategory(ctx, cat)
	if err != nil {
		return Category{}, err
	}

	return created, nil
}

func (s *Service) GetCategory(ctx context.Context, id string) (Category, error) {
	if id == "" {
		return Category{}, ErrCategoryNotFound
	}
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *Service) ListCategories(ctx context.Context, activeOnly bool) ([]Category, error) {
	return s.repo.ListCategories(ctx, activeOnly)
}

func (s *Service) UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (Category, error) {
	if id == "" {
		return Category{}, ErrCategoryNotFound
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return Category{}, ErrInvalidName
		}
		if utf8.RuneCountInString(name) > 100 {
			return Category{}, fmt.Errorf("name exceeds maximum length of 100 characters")
		}
		req.Name = &name
	}

	if req.Slug != nil {
		slug := strings.ToLower(strings.TrimSpace(*req.Slug))
		if slug == "" {
			return Category{}, ErrInvalidSlug
		}
		if !ValidateSlug(slug) {
			return Category{}, fmt.Errorf("invalid slug format")
		}
		req.Slug = &slug
	}

	cat, err := s.repo.UpdateCategory(ctx, id, req)
	if err != nil {
		return Category{}, err
	}

	return cat, nil
}

func (s *Service) DeactivateCategory(ctx context.Context, id string) error {
	if id == "" {
		return ErrCategoryNotFound
	}

	count, err := s.repo.CountActiveDesignsInCategory(ctx, id)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrCategoryHasDesigns
	}

	return s.repo.DeactivateCategory(ctx, id)
}
