package database

import (
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect creates a database connection and runs migrations
func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DB_URL))
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := db.AutoMigrate(&ImageModel{}); err != nil {
		return nil, err
	}

	return db, nil
}

