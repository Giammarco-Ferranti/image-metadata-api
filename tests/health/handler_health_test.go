package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/health"
)


func TestHealthCheckHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	
	health.HandlerHealth(w, request)

	resp := w.Result()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	} 

}
