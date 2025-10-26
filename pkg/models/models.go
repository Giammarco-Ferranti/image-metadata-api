package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)


type ImageProcess struct {
	ID 					uuid.UUID `gorm:"primaryKey"`
	CreatedAt		time.Time `gorm:"not null"`
	UpdatedAt 	time.Time `gorm:"not null"`
	Url					string
	Status 			string 		`gorm:"check:status IN ('done', 'failed', 'pending', 'in process')"`
	Width 			sql.NullInt16
	Height			sql.NullInt16
	Format			sql.NullString
}

type ImageProcessResponse struct {
	ID 					uuid.UUID `json:"id"`
	CreatedAt		string 		`json:"created_at"`
	UpdatedAt 	string 		`json:"updated_at"`
	Url					string		`json:"url"`
	Status 			string 		`json:"status"`
	Width 			*int16		`json:"width,omitempty"`
	Height			*int16		`json:"height,omitempty"`
	Format			*string		`json:"format,omitempty"`
}

// ToResponse converts ImageProcess to ImageProcessResponse
func (img *ImageProcess) ToResponse() ImageProcessResponse {
	response := ImageProcessResponse{
		ID: img.ID,
		CreatedAt: img.CreatedAt.Format(time.RFC3339),
		UpdatedAt: img.UpdatedAt.Format(time.RFC3339),
		Url: img.Url,
		Status: img.Status,
	}

	if img.Width.Valid {
		val := img.Width.Int16
		response.Width = &val
	}

	if img.Height.Valid {
		val := img.Height.Int16
		response.Height = &val
	}

	if img.Format.Valid {
		val := img.Format.String
		response.Format = &val
	}


	return response

}

// MarshalJSON provides custom JSON marshaling for ImageProcess
func (img ImageProcess) MarshalJSON() ([]byte, error) {
	response := img.ToResponse()
	return json.Marshal(response)
}