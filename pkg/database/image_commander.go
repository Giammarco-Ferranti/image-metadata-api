package database

import (
	"context"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
	"github.com/google/uuid"
)

type imageCommander struct {
	store Store
}

// NewImageCommander creates a new ImageCommander
func NewImageCommander(store Store) domain.ImageCommander {
	return &imageCommander{store: store}
}

// CreateImage implements domain.ImageCommander
func (c *imageCommander) CreateImage(ctx context.Context, url string) (*domain.Image, error) {
	now := time.Now().UTC()
	image := &domain.Image{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Url:       url,
		Status:    "pending",
		Width:     nil,
		Height:    nil,
		Format:    nil,
	}

	// Validate the image
	if err := image.Validate(); err != nil {
		return nil, err
	}

	// Save using repository within a transaction
	err := c.store.Atomic(func(s Store) error {
		return s.ImageRepository().Create(ctx, image)
	})

	if err != nil {
		return nil, err
	}

	return image, nil
}

