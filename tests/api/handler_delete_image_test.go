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
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandlerDeleteImageFails(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&models.ImageProcess{})

	imageMock := models.ImageProcess{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url: "http://example.com",
		Status: "in process",
		Width: sql.NullInt16{
			Int16: 200,
			Valid: true,
		},
		Height: sql.NullInt16{
			Int16: 200,
			Valid: true,
		},
		Format: sql.NullString{
			String: "png",
			Valid: true,
		},
	}

	db.Create(&imageMock)

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/image/%v", imageMock.ID), nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", imageMock.ID.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := api.Handler{DB: db}
	

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
	db.AutoMigrate(&models.ImageProcess{})

	imageMock := models.ImageProcess{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url: "http://example.com",
		Status: "done",
	}

	db.Create(&imageMock)

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/image/%v", imageMock.ID), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", imageMock.ID.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := api.Handler{DB: db}
	handler.HandleDeleteImage(w, request)

	resp := w.Result()

	// Assert that the status code is 200 (Success)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}