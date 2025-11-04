package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
)

// Handler that add url to the database
func (h Handler) HandlerAddImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()


	var req CreateImageRequest
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		slog.Error("Error decoding request body", "error", err)
		responses.RespondWithError(w, 400, "Error decoding request body")
		return
	}

	image, err := h.Commander.CreateImage(ctx, req.URL)
	if err != nil {
		// Domain validation errors come from Validate() method
		slog.Error("Error creating image", "error", err)
		responses.RespondWithError(w, 400, err.Error())
		return
	}

	response := NewImageResponse(image)
	responses.RespondWithJson(w, 201, ItemResponse{Data: response})
}
