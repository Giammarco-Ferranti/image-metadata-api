package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/api"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/database"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandlerDeleteImageFails(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&database.ImageModel{})

	store := database.NewStore(db)
	repository := store.ImageRepository()
	commander := database.NewImageCommander(store)

	ctx := context.Background()
	width := int16(200)
	height := int16(200)
	format := "png"
	imageMock := &domain.Image{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       "http://example.com",
		Status:    "in process",
		Width:     &width,
		Height:    &height,
		Format:    &format,
	}

	repository.Create(ctx, imageMock)

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/image/%v", imageMock.ID), nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", imageMock.ID.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := api.Handler{
		Querier:   repository,
		Commander: commander,
	}

	handler.HandleDeleteImage(w, request)

	resp := w.Result()

	bodyBytes, _ := io.ReadAll(resp.Body)

	// The error response has the structure: {"error": "message"}
	var errorResponse struct {
		Error string `json:"error"`
	}
	json.Unmarshal(bodyBytes, &errorResponse)

	// Assert that the status code is 400 (Bad Request)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	// Assert that the error message is correct
	expectedError := "Cannot delete image with status 'in process'"
	if errorResponse.Error != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, errorResponse.Error)
	}

}

func TestHandlerDeleteImageSuccess(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&database.ImageModel{})

	store := database.NewStore(db)
	repository := store.ImageRepository()
	commander := database.NewImageCommander(store)

	ctx := context.Background()
	imageMock := &domain.Image{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       "http://example.com",
		Status:    "done",
		Width:     nil,
		Height:    nil,
		Format:    nil,
	}

	repository.Create(ctx, imageMock)

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/image/%v", imageMock.ID), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", imageMock.ID.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := api.Handler{
		Querier:   repository,
		Commander: commander,
	}
	handler.HandleDeleteImage(w, request)

	resp := w.Result()

	// Assert that the status code is 200 (Success)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}