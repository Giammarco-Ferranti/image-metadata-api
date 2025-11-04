package health

import (
	"net/http"

	"gorm.io/gorm"
)

type HealthHandler struct {
	DB *gorm.DB
}

// func HandlerHealth(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(200)
// 	w.Write([]byte("OK"))
// }

func (h *HealthHandler) HandlerHealth(w http.ResponseWriter, r *http.Request) {
	sqlDb, err := h.DB.DB()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("KO"))
		return
	}

	err = sqlDb.Ping()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("KO"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}