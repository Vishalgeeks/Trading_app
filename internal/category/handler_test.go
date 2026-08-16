package category

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler_CreateCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	body := CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/categories", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCategory(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp CategoryResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, "Bridal", resp.Name)
}

func TestHandler_GetCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	cat, _ := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/categories/"+cat.ID, nil)
	w := httptest.NewRecorder()

	h.GetCategory(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_ListCategories(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	_, _ = svc.CreateCategory(context.Background(), CreateCategoryRequest{Name: "A", Slug: "a"})
	_, _ = svc.CreateCategory(context.Background(), CreateCategoryRequest{Name: "B", Slug: "b"})

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	w := httptest.NewRecorder()

	h.ListCategories(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err)
	require.GreaterOrEqual(t, float64(2), result["count"].(float64))
}

func TestHandler_UpdateCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	cat, _ := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	})

	desc := "Updated desc"
	body := UpdateCategoryRequest{Description: &desc}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/api/admin/categories/"+cat.ID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateCategory(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_DeleteCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	cat, _ := svc.CreateCategory(context.Background(), CreateCategoryRequest{
		Name: "Bridal",
		Slug: "bridal",
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/categories/"+cat.ID, nil)
	w := httptest.NewRecorder()

	h.DeleteCategory(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
