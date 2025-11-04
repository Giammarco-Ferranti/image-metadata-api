package responses

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		slog.Error("Responding with 5xx error", "message", msg)
	}

	type ErrorMsg struct {
		Error string `json:"error"`
	}

	RespondWithJson(w, code, ErrorMsg{Error: msg})
}


func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		slog.Error("Failed to marshal JSON response", "error", err, "payload", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

