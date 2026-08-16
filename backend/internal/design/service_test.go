package design

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"mehndi-booking-backend/internal/category"
)

type mockCategoryRepo struct {
	categories map[string]category.Category
}

func newMockCategoryRepo() *mockCategoryRepo {
	return &mockCategoryRepo{categories: make(map[string]category.Category)}
}

func (m *mockCategoryRepo) CreateCategory(ctx context.Context, cat category.Category) (category.Category, error) {
	for _, c := range m.categories {
		if c.Name == cat.Name {
			return category.Category{}, category.ErrDuplicateName
		}
		if c.Slug == cat.Slug {
			return category.Category{}, category.ErrDuplicateSlug
		}
	}
	cat.ID = "mock-cat-" + cat.Slug
	m.categories[cat.ID] = cat
	return cat, nil
}

func (m *mockCategoryRepo) GetCategory(ctx context.Context, id string) (category.Category, error) {
	cat, ok := m.categories[id]
	if !ok {
		return category.Category{}, category.ErrCategoryNotFound
	}
	return cat, nil
}

func (m *mockCategoryRepo) GetCategoryByID(ctx context.Context, id string) (category.Category, error) {
	return m.GetCategory(ctx, id)
}

func (m *mockCategoryRepo) GetCategoryBySlug(ctx context.Context, slug string) (category.Category, error) {
	for _, cat := range m.categories {
		if cat.Slug == slug {
			return cat, nil
		}
	}
	return category.Category{}, category.ErrCategoryNotFound
}

func (m *mockCategoryRepo) ListCategories(ctx context.Context, activeOnly bool) ([]category.Category, error) {
	var cats []category.Category
	for _, cat := range m.categories {
		if activeOnly && !cat.IsActive {
			continue
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func (m *mockCategoryRepo) UpdateCategory(ctx context.Context, id string, req category.UpdateCategoryRequest) (category.Category, error) {
	cat, ok := m.categories[id]
	if !ok {
		return category.Category{}, category.ErrCategoryNotFound
	}
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Slug != nil {
		cat.Slug = *req.Slug
	}
	if req.Description != nil {
		cat.Description = req.Description
	}
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}
	m.categories[id] = cat
	return cat, nil
}

func (m *mockCategoryRepo) DeactivateCategory(ctx context.Context, id string) error {
	cat, ok := m.categories[id]
	if !ok {
		return category.ErrCategoryNotFound
	}
	cat.IsActive = false
	m.categories[id] = cat
	return nil
}

func (m *mockCategoryRepo) CountActiveDesignsInCategory(ctx context.Context, categoryID string) (int, error) {
	return 0, nil
}

type mockDesignRepo struct {
	designs map[string]Design
}

func newMockDesignRepo() *mockDesignRepo {
	return &mockDesignRepo{designs: make(map[string]Design)}
}

func (m *mockDesignRepo) CreateDesign(ctx context.Context, des Design) (Design, error) {
	for _, d := range m.designs {
		if d.Slug == des.Slug {
			return Design{}, ErrDuplicateSlug
		}
	}
	des.ID = "mock-design-" + des.Slug
	m.designs[des.ID] = des
	return des, nil
}

func (m *mockDesignRepo) GetDesignByID(ctx context.Context, id string) (Design, error) {
	des, ok := m.designs[id]
	if !ok {
		return Design{}, ErrDesignNotFound
	}
	return des, nil
}

func (m *mockDesignRepo) GetDesignBySlug(ctx context.Context, slug string) (Design, error) {
	for _, des := range m.designs {
		if des.Slug == slug {
			return des, nil
		}
	}
	return Design{}, ErrDesignNotFound
}

func (m *mockDesignRepo) ListDesigns(ctx context.Context, categoryID *string, searchQuery *string, page, limit int) ([]Design, error) {
	var designs []Design
	for _, des := range m.designs {
		if !des.IsActive {
			continue
		}
		if categoryID != nil && des.CategoryID != *categoryID {
			continue
		}
		if searchQuery != nil && *searchQuery != "" {
			q := strings.ToLower(*searchQuery)
			found := strings.Contains(strings.ToLower(des.Name), q)
			if !found && des.Description != nil {
				found = strings.Contains(strings.ToLower(*des.Description), q)
			}
			if !found {
				continue
			}
		}
		designs = append(designs, des)
	}

	if limit > 0 {
		start := (page - 1) * limit
		if start > len(designs) {
			return []Design{}, nil
		}
		end := start + limit
		if end > len(designs) {
			end = len(designs)
		}
		designs = designs[start:end]
	}

	return designs, nil
}

func (m *mockDesignRepo) ListActiveDesigns(ctx context.Context) ([]Design, error) {
	return m.ListDesigns(ctx, nil, nil, 0, 0)
}

func (m *mockDesignRepo) ListDesignsByCategory(ctx context.Context, categoryID string) ([]Design, error) {
	return m.ListDesigns(ctx, &categoryID, nil, 0, 0)
}

func (m *mockDesignRepo) SearchDesigns(ctx context.Context, query string) ([]Design, error) {
	return m.ListDesigns(ctx, nil, &query, 0, 0)
}

func (m *mockDesignRepo) UpdateDesign(ctx context.Context, id string, req UpdateDesignRequest) (Design, error) {
	des, ok := m.designs[id]
	if !ok {
		return Design{}, ErrDesignNotFound
	}
	if req.CategoryID != nil {
		des.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		des.Name = *req.Name
	}
	if req.Slug != nil {
		des.Slug = *req.Slug
	}
	if req.Description != nil {
		des.Description = req.Description
	}
	if req.ImageURL != nil {
		des.ImageURL = *req.ImageURL
	}
	if req.Price != nil {
		des.Price = *req.Price
	}
	if req.DurationMinutes != nil {
		des.DurationMinutes = *req.DurationMinutes
	}
	if req.IsActive != nil {
		des.IsActive = *req.IsActive
	}
	m.designs[id] = des
	return des, nil
}

func (m *mockDesignRepo) DeactivateDesign(ctx context.Context, id string) error {
	des, ok := m.designs[id]
	if !ok {
		return ErrDesignNotFound
	}
	des.IsActive = false
	m.designs[id] = des
	return nil
}

func (m *mockDesignRepo) GetCategory(ctx context.Context, id string) (category.Category, error) {
	return category.Category{}, category.ErrCategoryNotFound
}

func TestService_CreateDesign(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})

	des, err := svc.CreateDesign(context.Background(), CreateDesignRequest{
		CategoryID:      cat.ID,
		Name:            "Bridal Full Hand",
		Slug:            "bridal-full-hand",
		ImageURL:        "https://example.com/design.jpg",
		Price:           "2500.00",
		DurationMinutes: 180,
	})
	require.NoError(t, err)
	require.Equal(t, "Bridal Full Hand", des.Name)
	require.Equal(t, "bridal-full-hand", des.Slug)
}

func TestService_CreateDesign_InvalidCategory(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	_, err := svc.CreateDesign(context.Background(), CreateDesignRequest{
		CategoryID:      "invalid-uuid",
		Name:            "Test Design",
		Slug:            "test-design",
		ImageURL:        "https://example.com/design.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
	})
	require.ErrorIs(t, err, ErrCategoryNotFound)
}

func TestService_CreateDesign_InactiveCategory(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Inactive",
		Slug:     "inactive",
		IsActive: false,
	})

	_, err := svc.CreateDesign(context.Background(), CreateDesignRequest{
		CategoryID:      cat.ID,
		Name:            "Test Design",
		Slug:            "test-design",
		ImageURL:        "https://example.com/design.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
	})
	require.ErrorIs(t, err, ErrCategoryInactive)
}

func TestService_CreateDesign_DuplicateSlug(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})

	_, _ = svc.CreateDesign(context.Background(), CreateDesignRequest{
		CategoryID:      cat.ID,
		Name:            "Design 1",
		Slug:            "design-1",
		ImageURL:        "https://example.com/design1.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
	})

	_, err := svc.CreateDesign(context.Background(), CreateDesignRequest{
		CategoryID:      cat.ID,
		Name:            "Design 2",
		Slug:            "design-1",
		ImageURL:        "https://example.com/design2.jpg",
		Price:           "2000.00",
		DurationMinutes: 90,
	})
	require.ErrorIs(t, err, ErrDuplicateSlug)
}

func TestService_CreateDesign_InvalidPrice(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})

	_, err := svc.CreateDesign(context.Background(), CreateDesignRequest{
		CategoryID:      cat.ID,
		Name:            "Test Design",
		Slug:            "test-design",
		ImageURL:        "https://example.com/design.jpg",
		Price:           "-100",
		DurationMinutes: 60,
	})
	require.Error(t, err)
}

func TestService_CreateDesign_InvalidDuration(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})

	_, err := svc.CreateDesign(context.Background(), CreateDesignRequest{
		CategoryID:      cat.ID,
		Name:            "Test Design",
		Slug:            "test-design",
		ImageURL:        "https://example.com/design.jpg",
		Price:           "1000.00",
		DurationMinutes: 0,
	})
	require.Error(t, err)
}

func TestService_GetDesign(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})
	created, _ := desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat.ID,
		Name:            "Bridal Full Hand",
		Slug:            "bridal-full-hand",
		ImageURL:        "https://example.com/design.jpg",
		Price:           "2500.00",
		DurationMinutes: 180,
		IsActive:        true,
	})

	found, err := svc.GetDesign(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, "Bridal Full Hand", found.Name)
}

func TestService_ListDesigns(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})

	_, _ = desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat.ID,
		Name:            "Design 1",
		Slug:            "design-1",
		ImageURL:        "https://example.com/design1.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})
	_, _ = desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat.ID,
		Name:            "Design 2",
		Slug:            "design-2",
		ImageURL:        "https://example.com/design2.jpg",
		Price:           "2000.00",
		DurationMinutes: 90,
		IsActive:        true,
	})

	designs, err := svc.ListDesigns(context.Background(), nil, nil, 1, 20)
	require.NoError(t, err)
	require.Len(t, designs, 2)
}

func TestService_ListDesigns_FilterByCategory(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat1, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})
	cat2, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Arabic",
		Slug:     "arabic",
		IsActive: true,
	})

	_, _ = desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat1.ID,
		Name:            "Bridal Design",
		Slug:            "bridal-design",
		ImageURL:        "https://example.com/bridal.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})
	_, _ = desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat2.ID,
		Name:            "Arabic Design",
		Slug:            "arabic-design",
		ImageURL:        "https://example.com/arabic.jpg",
		Price:           "2000.00",
		DurationMinutes: 90,
		IsActive:        true,
	})

	designs, err := svc.ListDesigns(context.Background(), &cat1.ID, nil, 1, 20)
	require.NoError(t, err)
	require.Len(t, designs, 1)
}

func TestService_SearchDesigns(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})

	_, _ = desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat.ID,
		Name:            "Bridal Full Hand",
		Slug:            "bridal-full-hand",
		ImageURL:        "https://example.com/bridal.jpg",
		Price:           "2500.00",
		DurationMinutes: 180,
		IsActive:        true,
	})
	_, _ = desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat.ID,
		Name:            "Arabic Light",
		Slug:            "arabic-light",
		ImageURL:        "https://example.com/arabic.jpg",
		Price:           "1500.00",
		DurationMinutes: 120,
		IsActive:        true,
	})

	designs, err := svc.SearchDesigns(context.Background(), "bridal")
	require.NoError(t, err)
	require.Len(t, designs, 1)
}

func TestService_UpdateDesign(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})
	created, _ := desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat.ID,
		Name:            "Original Design",
		Slug:            "original-design",
		ImageURL:        "https://example.com/original.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})

	newPrice := "1500.00"
	updated, err := svc.UpdateDesign(context.Background(), created.ID, UpdateDesignRequest{
		Price: &newPrice,
	})
	require.NoError(t, err)
	require.Equal(t, "1500.00", updated.Price)
}

func TestService_DeactivateDesign(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})
	created, _ := desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat.ID,
		Name:            "Test Design",
		Slug:            "test-design",
		ImageURL:        "https://example.com/test.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})

	err := svc.DeactivateDesign(context.Background(), created.ID)
	require.NoError(t, err)

	found, _ := desRepo.GetDesignByID(context.Background(), created.ID)
	require.False(t, found.IsActive)
}

func TestService_ListDesigns_Pagination(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})

	for i := 0; i < 5; i++ {
		slug := fmt.Sprintf("design-%d", i)
		desRepo.CreateDesign(context.Background(), Design{
			CategoryID:      cat.ID,
			Name:            "Design " + string(rune('0'+i)),
			Slug:            slug,
			ImageURL:        "https://example.com/" + slug + ".jpg",
			Price:           "1000.00",
			DurationMinutes: 60,
			IsActive:        true,
		})
	}

	designs, err := svc.ListDesigns(context.Background(), nil, nil, 1, 3)
	require.NoError(t, err)
	require.Len(t, designs, 3)
}
