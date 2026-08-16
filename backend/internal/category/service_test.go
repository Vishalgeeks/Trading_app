package category

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockCategoryRepo struct {
	categories map[string]Category
}

func newMockCategoryRepo() *mockCategoryRepo {
	return &mockCategoryRepo{categories: make(map[string]Category)}
}

func (m *mockCategoryRepo) CreateCategory(ctx context.Context, cat Category) (Category, error) {
	for _, c := range m.categories {
		if c.Name == cat.Name {
			return Category{}, ErrDuplicateName
		}
		if c.Slug == cat.Slug {
			return Category{}, ErrDuplicateSlug
		}
	}
	cat.ID = "mock-cat-" + cat.Slug
	m.categories[cat.ID] = cat
	return cat, nil
}

func (m *mockCategoryRepo) GetCategoryByID(ctx context.Context, id string) (Category, error) {
	cat, ok := m.categories[id]
	if !ok {
		return Category{}, ErrCategoryNotFound
	}
	return cat, nil
}

func (m *mockCategoryRepo) GetCategoryBySlug(ctx context.Context, slug string) (Category, error) {
	for _, cat := range m.categories {
		if cat.Slug == slug {
			return cat, nil
		}
	}
	return Category{}, ErrCategoryNotFound
}

func (m *mockCategoryRepo) ListCategories(ctx context.Context, activeOnly bool) ([]Category, error) {
	var cats []Category
	for _, cat := range m.categories {
		if activeOnly && !cat.IsActive {
			continue
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func (m *mockCategoryRepo) UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (Category, error) {
	cat, ok := m.categories[id]
	if !ok {
		return Category{}, ErrCategoryNotFound
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
		return ErrCategoryNotFound
	}
	cat.IsActive = false
	m.categories[id] = cat
	return nil
}

func (m *mockCategoryRepo) CountActiveDesignsInCategory(ctx context.Context, categoryID string) (int, error) {
	return 0, nil
}

func TestService_CreateCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)

	cat, err := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	})
	require.NoError(t, err)
	require.Equal(t, "Bridal", cat.Name)
	require.Equal(t, "bridal", cat.Slug)
}

func TestService_CreateCategory_DuplicateName(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)

	_, err := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	})
	require.NoError(t, err)

	_, err = svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal-2",
	})
	require.ErrorIs(t, err, ErrDuplicateName)
}

func TestService_CreateCategory_InvalidSlug(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)

	_, err := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "INVALID SLUG",
	})
	require.Error(t, err)
}

func TestService_GetCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)

	cat, err := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	})
	require.NoError(t, err)

	found, err := svc.GetCategory(context.Background(), cat.ID)
	require.NoError(t, err)
	require.Equal(t, "Bridal", found.Name)
}

func TestService_ListCategories(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)

	cat1, err1 := svc.CreateCategory(context.Background(), CreateCategoryRequest{Name: "Alpha", Slug: "alpha"})
	t.Logf("cat1: %v, err1: %v", cat1, err1)

	cat2, err2 := svc.CreateCategory(context.Background(), CreateCategoryRequest{Name: "Beta", Slug: "beta"})
	t.Logf("cat2: %v, err2: %v", cat2, err2)

	cats, err := svc.ListCategories(context.Background(), false)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(cats), 2)
}

func TestService_UpdateCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)

	cat, err := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	})
	require.NoError(t, err)

	desc := "Updated desc"
	updated, err := svc.UpdateCategory(context.Background(), cat.ID, UpdateCategoryRequest{
		Description: &desc,
	})
	require.NoError(t, err)
	require.Equal(t, "Updated desc", *updated.Description)
}

func TestService_DeactivateCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)

	cat, err := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	})
	require.NoError(t, err)

	err = svc.DeactivateCategory(context.Background(), cat.ID)
	require.NoError(t, err)

	found, err := svc.GetCategory(context.Background(), cat.ID)
	require.NoError(t, err)
	require.False(t, found.IsActive)
}
