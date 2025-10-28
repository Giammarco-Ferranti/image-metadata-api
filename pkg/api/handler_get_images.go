package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
)


func (h Handler) HandlerGetImages (w http.ResponseWriter, r *http.Request) {

	totalStr := r.URL.Query().Get("total")
	offsetStr := r.URL.Query().Get("offset")


	total, err := strconv.Atoi(totalStr)
	if err != nil {
		total = 100 //default
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0 //default
	}
	if total > 100 {
		responses.RespondWithError(w, 400, "Total cannot be more that 100")
		return
	}


	var images []models.ImageProcess
	err = h.DB.Limit(total).Offset(offset).Find(&images).Error
	if err != nil {
		log.Println("Couldn't retrieve images: ", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Couldn't retrieve images: %v", err))
		return
	}
	
	// Convert to response format
	imageResponses := make([]models.ImageProcessResponse, len(images))
	for i, img := range images {
		imageResponses[i] = img.ToResponse()
	}
	
	responses.RespondWithJson(w, 200, responses.PaginatedResponse{
		Data: imageResponses,
		Meta: responses.PaginationMeta{
			Total: total,
			Offset: offset,
		},
	})
}
