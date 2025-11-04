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
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandlerGetImagesSuccess(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&database.ImageModel{})

	store := database.NewStore(db)
	repository := store.ImageRepository()
	commander := database.NewImageCommander(store)

	ctx := context.Background()
	// Create test images
	for i := 0; i < 3; i++ {
		width := int16(100)
		height := int16(100)
		format := "png"
		image := &domain.Image{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Url:       fmt.Sprintf("http://example.com/%d", i),
			Status:    "done",
			Width:     &width,
			Height:    &height,
			Format:    &format,
		}
		repository.Create(ctx, image)
	}

	request := httptest.NewRequest(http.MethodGet, "/images?total=2&offset=0", nil)
	w := httptest.NewRecorder()
	handler := api.Handler{
		Querier:   repository,
		Commander: commander,
	}
	handler.HandlerGetImages(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var response api.PaginatedImageResponse
	json.Unmarshal(bodyBytes, &response)

	if len(response.Data) != 2 {
		t.Errorf("expected 2 images, got %d", len(response.Data))
	}
}

func TestHandlerGetImagesInvalidTotal(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&database.ImageModel{})

	store := database.NewStore(db)
	repository := store.ImageRepository()
	commander := database.NewImageCommander(store)

	request := httptest.NewRequest(http.MethodGet, "/images?total=101&offset=0", nil)
	w := httptest.NewRecorder()
	handler := api.Handler{
		Querier:   repository,
		Commander: commander,
	}
	handler.HandlerGetImages(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

