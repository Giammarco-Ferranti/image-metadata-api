package domain

import (
	"context"

	"github.com/google/uuid"
)

type ImageQuerier interface {
	FindById (ctx context.Context, id uuid.UUID) (*Image, error)
	FindAll (ctx context.Context, total, offset int) ([]*Image, error)
	FindPendingImages(ctx context.Context, limit int) ([]*Image, error)
}

type ImageRepository interface {
	ImageQuerier
	Create (ctx context.Context, image *Image) error
	Save (ctx context.Context, image *Image) error
	Delete(ctx context.Context, id uuid.UUID) error
}