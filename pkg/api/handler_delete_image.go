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

func (h Handler) HandleDeleteImage(w http.ResponseWriter, r *http.Request) {
	
	imageIdString := chi.URLParam(r, "id")
	imageId, err := uuid.Parse(imageIdString)

	if err != nil {
		log.Println("Couldn't parse uuid:", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse id: %v", err))
		return
	}

	var image models.ImageProcess
	err = h.DB.Where("id = ?", imageId).Find(&image).Error

	if err != nil {
		log.Println("Couldn't get image:", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Couldn't get image: %v", err))
		return
	}

	if image.Status == "in process" {
		log.Println("Cannot delete image with status 'in process'")
		responses.RespondWithError(w, 400, "Cannot delete image with status 'in process'")
		return
	}

	err = h.DB.Where("id = ?", imageId).Delete(&image).Error

	if err != nil {
		log.Println("Error deleting image:", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Error deleting image: %v", err))
		return
	}

	responses.RespondWithJson(w, 200, "Successfully deleted image")

}