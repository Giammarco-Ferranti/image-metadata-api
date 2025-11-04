package database

import (
	"database/sql"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
	"github.com/google/uuid"
)


type ImageModel struct {
	ID 					uuid.UUID `gorm:"primaryKey"`
	CreatedAt		time.Time `gorm:"not null"`
	UpdatedAt 	time.Time `gorm:"not null"`
	Url					string
	Status 			string 		`gorm:"check:status IN ('done', 'failed', 'pending', 'in process')"`
	Width 			sql.NullInt16
	Height			sql.NullInt16
	Format			sql.NullString
}

// ToDomain converts ImageModel to domain.Image
func (m *ImageModel) ToDomain() *domain.Image {
	domainImg := &domain.Image{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Url:       m.Url,
		Status:    m.Status,
	}

	// Convert nullable fields to pointers
	if m.Width.Valid {
		val := m.Width.Int16
		domainImg.Width = &val
	}

	if m.Height.Valid {
		val := m.Height.Int16
		domainImg.Height = &val
	}

	if m.Format.Valid {
		val := m.Format.String
		domainImg.Format = &val
	}

	return domainImg
}

// FromDomain converts domain.Image to ImageModel
func (m *ImageModel) FromDomain(img *domain.Image) {
	m.ID = img.ID
	m.CreatedAt = img.CreatedAt
	m.UpdatedAt = img.UpdatedAt
	m.Url = img.Url
	m.Status = img.Status

	// Convert pointers to nullable fields
	if img.Width != nil {
		m.Width = sql.NullInt16{Int16: *img.Width, Valid: true}
	} else {
		m.Width = sql.NullInt16{Valid: false}
	}

	if img.Height != nil {
		m.Height = sql.NullInt16{Int16: *img.Height, Valid: true}
	} else {
		m.Height = sql.NullInt16{Valid: false}
	}

	if img.Format != nil {
		m.Format = sql.NullString{String: *img.Format, Valid: true}
	} else {
		m.Format = sql.NullString{Valid: false}
	}
}