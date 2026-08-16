package design

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"mehndi-booking-backend/internal/category"
)

func TestHandler_CreateDesign(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)
	h := NewHandler(svc)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})

	body := CreateDesignRequest{
		CategoryID:      cat.ID,
		Name:            "Bridal Full Hand",
		Slug:            "bridal-full-hand",
		ImageURL:        "https://example.com/design.jpg",
		Price:           "2500.00",
		DurationMinutes: 180,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/designs", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateDesign(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp DesignResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "Bridal Full Hand", resp.Name)
}

func TestHandler_GetDesign(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)
	h := NewHandler(svc)

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

	req := httptest.NewRequest(http.MethodGet, "/api/designs/"+created.ID, nil)
	w := httptest.NewRecorder()

	h.GetDesign(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_ListDesigns(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)
	h := NewHandler(svc)

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

	req := httptest.NewRequest(http.MethodGet, "/api/designs?page=1&limit=20", nil)
	w := httptest.NewRecorder()

	h.ListDesigns(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	require.GreaterOrEqual(t, float64(1), result["count"].(float64))
}

func TestHandler_SearchDesigns(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)
	h := NewHandler(svc)

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

	req := httptest.NewRequest(http.MethodGet, "/api/designs/search?q=bridal", nil)
	w := httptest.NewRecorder()

	h.SearchDesigns(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	require.GreaterOrEqual(t, float64(1), result["count"].(float64))
}

func TestHandler_ListDesignsByCategory(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)
	h := NewHandler(svc)

	cat, _ := catRepo.CreateCategory(context.Background(), category.Category{
		Name:     "Bridal",
		Slug:     "bridal",
		IsActive: true,
	})
	_, _ = desRepo.CreateDesign(context.Background(), Design{
		CategoryID:      cat.ID,
		Name:            "Bridal Design",
		Slug:            "bridal-design",
		ImageURL:        "https://example.com/bridal.jpg",
		Price:           "1000.00",
		DurationMinutes: 60,
		IsActive:        true,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/designs/category/"+cat.ID, nil)
	w := httptest.NewRecorder()

	h.ListDesignsByCategory(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	require.GreaterOrEqual(t, float64(1), result["count"].(float64))
}

func TestHandler_AdminUpdateDesign(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)
	h := NewHandler(svc)

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
	body := UpdateDesignRequest{Price: &newPrice}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/api/admin/designs/"+created.ID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.AdminUpdateDesign(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp DesignResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "1500.00", resp.Price)
}

func TestHandler_AdminDeleteDesign(t *testing.T) {
	catRepo := newMockCategoryRepo()
	desRepo := newMockDesignRepo()
	svc := NewService(desRepo, catRepo)
	h := NewHandler(svc)

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

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/designs/"+created.ID, nil)
	w := httptest.NewRecorder()

	h.AdminDeleteDesign(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
