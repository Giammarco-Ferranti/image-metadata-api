package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/api"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandlerGetImageSuccess(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&models.ImageProcess{})

	imageMock := models.ImageProcess{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url: "http://example.com",
		Status: "done",
		Width: sql.NullInt16{Int16: 100, Valid: true},
		Height: sql.NullInt16{Int16: 100, Valid: true},
		Format: sql.NullString{String: "png", Valid: true},
	}
	db.Create(&imageMock)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/image/%v", imageMock.ID), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", imageMock.ID.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := api.Handler{DB: db}
	handler.HandlerGetImage(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var response responses.ItemResponse
	json.Unmarshal(bodyBytes, &response)

	if response.Data.ID != imageMock.ID {
		t.Errorf("expected id %v, got %v", imageMock.ID, response.Data.ID)
	}
}

func TestHandlerGetImageNotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&models.ImageProcess{})

	fakeID := uuid.New()
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/image/%v", fakeID), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", fakeID.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := api.Handler{DB: db}
	handler.HandlerGetImage(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

