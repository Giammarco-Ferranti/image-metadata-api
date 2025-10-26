package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
)


func (h Handler) HandlerGetImages (w http.ResponseWriter, r *http.Request) {
	var images []models.ImageProcess
	err := h.DB.Find(&images).Limit(100).Error
	if err != nil {
		log.Println("Couldn't retrieve images: ", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Couldn't retrieve images: %v", err))
		return
	}
	
	responses.RespondWithJson(w, 200, images)
}
