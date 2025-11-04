package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/api"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandlerAddImage(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&database.ImageModel{})

	store := database.NewStore(db)
	repository := store.ImageRepository()
	commander := database.NewImageCommander(store)

	body := map[string]string{
		"url": "http://example.com",
	}

	jsonBody, _ := json.Marshal(body)

	request := httptest.NewRequest(http.MethodPost, "/image", bytes.NewReader(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler := api.Handler{
		Querier:   repository,
		Commander: commander,
	}

	handler.HandlerAddImage(w, request)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	var response api.ItemResponse
	json.Unmarshal(bodyBytes, &response)

	if response.Data.Url != body["url"] {
		t.Errorf("expected url %s, got %s", body["url"], response.Data.Url)
	}

}