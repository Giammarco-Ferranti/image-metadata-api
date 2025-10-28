package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (h *Handler) HandlerGetImage(w http.ResponseWriter, r *http.Request) {
	imageIdString := chi.URLParam(r, "id")
	imageId, err := uuid.Parse(imageIdString)

	if err != nil {
		log.Println("Invalid image id: ", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Invalid image id, %v", err))
		return
	}

	image := models.ImageProcess{}

	err = h.DB.Where("id = ?", imageId).First(&image).Error

	if err != nil {
		log.Println("Error retrieving image: ", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Error retrieving image: %v", err))
		return
	}

	responses.RespondWithJson(w, 200, responses.ItemResponse{
		Data: image.ToResponse(),
	})
	
}