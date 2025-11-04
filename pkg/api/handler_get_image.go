package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (h *Handler) HandlerGetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	imageIdString := chi.URLParam(r, "id")
	imageId, err := uuid.Parse(imageIdString)

	if err != nil {
		slog.Error("Invalid image id", "error", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Invalid image id, %v", err))
		return
	}


	image, err := h.Querier.FindById(ctx, imageId)

	if err != nil {
		slog.Error("Error retrieving image", "error", err)
		responses.RespondWithError(w, 500, fmt.Sprintf("Error retrieving image: %v", err))
		return
	}

	// Check if image was not found (returns nil, nil)
	if image == nil {
		responses.RespondWithError(w, 404, "Image not found")
		return
	}

	responses.RespondWithJson(w, 200, ItemResponse{
		Data: NewImageResponse(image),
	})
	
}