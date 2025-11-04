package health

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/health"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error creating test database: %v", err)
	}

	// Create health handler with test database
	healthHandler := &health.HealthHandler{DB: db}

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	healthHandler.HandlerHealth(w, request)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	if string(bodyBytes) != "OK" {
		t.Errorf("expected body 'OK', got %s", string(bodyBytes))
	}
}
