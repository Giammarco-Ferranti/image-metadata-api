package database

import (
	"context"
	"errors"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type imageRepository struct {
	DB *gorm.DB
}

// FindById implements ImageQuerier interface
func (r *imageRepository) FindById(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	var model ImageModel
	result := r.DB.WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

// FindAll implements ImageQuerier interface
func (r *imageRepository) FindAll(ctx context.Context, total, offset int) ([]*domain.Image, error) {
	var models []ImageModel
	result := r.DB.WithContext(ctx).
		Limit(total).
		Offset(offset).
		Find(&models)
	
	if result.Error != nil {
		return nil, result.Error
	}

	images := make([]*domain.Image, len(models))
	for i := range models {
		images[i] = models[i].ToDomain()
	}
	return images, nil
}

// FindPendingImages implements ImageQuerier interface
func (r *imageRepository) FindPendingImages(ctx context.Context, limit int) ([]*domain.Image, error) {
	var models []ImageModel
	result := r.DB.WithContext(ctx).
		Where("status = ?", "pending").
		Order("created_at asc").
		Limit(limit).
		Find(&models)
	
	if result.Error != nil {
		return nil, result.Error
	}

	images := make([]*domain.Image, len(models))
	for i := range models {
		images[i] = models[i].ToDomain()
	}
	return images, nil
}

// Create implements ImageRepository interface
func (r *imageRepository) Create(ctx context.Context, image *domain.Image) error {
	var model ImageModel
	model.FromDomain(image)
	
	result := r.DB.WithContext(ctx).Create(&model)
	return result.Error
}

// Save implements ImageRepository interface
func (r *imageRepository) Save(ctx context.Context, image *domain.Image) error {
	var model ImageModel
	model.FromDomain(image)
	
	result := r.DB.WithContext(ctx).Save(&model)
	return result.Error
}

// Delete implements ImageRepository interface
func (r *imageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&ImageModel{}, id)
	return result.Error
}

