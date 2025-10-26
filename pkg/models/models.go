package models

import (
	"database/sql"
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