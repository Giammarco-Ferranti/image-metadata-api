package api

import (
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
	"github.com/google/uuid"
)

type ImageResponse struct {
	ID 					uuid.UUID `json:"id"`
	CreatedAt		string 		`json:"created_at"`
	UpdatedAt 	string 		`json:"updated_at"`
	Url					string		`json:"url"`
	Status 			string 		`json:"status"`
	Width 			*int16		`json:"width,omitempty"`
	Height			*int16		`json:"height,omitempty"`
	Format			*string		`json:"format,omitempty"`
}

// NewImageResponse converts a domain.Image to ImageResponse
func NewImageResponse(img *domain.Image) ImageResponse {
	return ImageResponse{
		ID:        img.ID,
		CreatedAt: img.CreatedAt.Format(time.RFC3339),
		UpdatedAt: img.UpdatedAt.Format(time.RFC3339),
		Url:       img.Url,
		Status:    img.Status,
		Width:     img.Width,
		Height:    img.Height,
		Format:    img.Format,
	}
}

// PaginatedImageResponse is a paginated response for images
type PaginatedImageResponse struct {
	Data []ImageResponse `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Total        int `json:"total"`
	Offset       int `json:"offset"`
	ResultsCount int `json:"results_count"`
}

// ItemResponse is for single item responses
type ItemResponse struct {
	Data ImageResponse `json:"data"`
}

