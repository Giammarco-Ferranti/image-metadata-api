package api

import (
	"fmt"
	"log"
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
		log.Println("Invalid image id: ", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Invalid image id, %v", err))
		return
	}


	image, err := h.Querier.FindById(ctx, imageId)

	if err != nil {
		log.Println("Error retrieving image: ", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Error retrieving image: %v", err))
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