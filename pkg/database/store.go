package database

import (
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
	"gorm.io/gorm"
)

type store struct {
	DB *gorm.DB
}

type Store interface {
	Atomic(fn func(Store) error) error
	ImageRepository() domain.ImageRepository
}

func NewStore(db *gorm.DB) Store {
	return &store{DB: db}
}

// Atomic wraps the provided function in a GORM transaction
func (s *store) Atomic(fn func(Store) error) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		txStore := &store{DB: tx}
		return fn(txStore)
	})
}

// ImageRepository returns an ImageRepository instance
func (s *store) ImageRepository() domain.ImageRepository {
	return &imageRepository{DB: s.DB}
}