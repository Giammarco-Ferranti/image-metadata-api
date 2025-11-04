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

func TestHandlerGetImageSuccess(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&database.ImageModel{})

	store := database.NewStore(db)
	repository := store.ImageRepository()
	commander := database.NewImageCommander(store)

	ctx := context.Background()
	width := int16(100)
	height := int16(100)
	format := "png"
	imageMock := &domain.Image{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       "http://example.com",
		Status:    "done",
		Width:     &width,
		Height:    &height,
		Format:    &format,
	}
	repository.Create(ctx, imageMock)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/image/%v", imageMock.ID), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", imageMock.ID.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := api.Handler{
		Querier:   repository,
		Commander: commander,
	}
	handler.HandlerGetImage(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var response api.ItemResponse
	json.Unmarshal(bodyBytes, &response)

	if response.Data.ID != imageMock.ID {
		t.Errorf("expected id %v, got %v", imageMock.ID, response.Data.ID)
	}
}

func TestHandlerGetImageNotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&database.ImageModel{})

	store := database.NewStore(db)
	repository := store.ImageRepository()
	commander := database.NewImageCommander(store)

	fakeID := uuid.New()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/image/%v", fakeID), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", fakeID.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := api.Handler{
		Querier:   repository,
		Commander: commander,
	}
	handler.HandlerGetImage(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

