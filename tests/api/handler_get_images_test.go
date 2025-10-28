package api

import (
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
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandlerGetImagesSuccess(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&models.ImageProcess{})

	// Create test images
	for i := 0; i < 3; i++ {
		image := models.ImageProcess{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Url: fmt.Sprintf("http://example.com/%d", i),
			Status: "done",
			Width: sql.NullInt16{Int16: 100, Valid: true},
			Height: sql.NullInt16{Int16: 100, Valid: true},
			Format: sql.NullString{String: "png", Valid: true},
		}
		db.Create(&image)
	}

	request := httptest.NewRequest(http.MethodGet, "/images?total=2&offset=0", nil)
	w := httptest.NewRecorder()
	handler := api.Handler{DB: db}
	handler.HandlerGetImages(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var response responses.PaginatedResponse
	json.Unmarshal(bodyBytes, &response)

	if len(response.Data) != 2 {
		t.Errorf("expected 2 images, got %d", len(response.Data))
	}
}

func TestHandlerGetImagesInvalidTotal(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&models.ImageProcess{})

	request := httptest.NewRequest(http.MethodGet, "/images?total=101&offset=0", nil)
	w := httptest.NewRecorder()
	handler := api.Handler{DB: db}
	handler.HandlerGetImages(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

