package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
)


func (h Handler) HandlerGetImages (w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
	if total > 100 || total < 0 || offset < 0 {
		responses.RespondWithError(w, 400, "Total or offset not correct")
		return
	}

	
	images, err := h.Querier.FindAll(ctx, total, offset)
	
	if err != nil {
		log.Println("Couldn't retrieve images: ", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Couldn't retrieve images: %v", err))
		return
	}
	
	// Convert to response format
	imageResponses := make([]ImageResponse, len(images))
	for i, img := range images {
		imageResponses[i] = NewImageResponse(img)
	}
	
	responses.RespondWithJson(w, 200, PaginatedImageResponse{
		Data: imageResponses,
		Meta: PaginationMeta{
			Total:        total,
			Offset:       offset,
			ResultsCount: len(images),
		},
	})
}
