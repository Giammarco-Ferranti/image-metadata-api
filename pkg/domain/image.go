package domain

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
)


type Image struct {
	ID 					uuid.UUID 
	CreatedAt		time.Time 
	UpdatedAt 	time.Time 
	Url					string
	Status 			string 		
	Width 			*int16
	Height			*int16
	Format			*string
}

func (img Image) Validate() error {
	// Validate URL is not empty
	if len(img.Url) == 0 {
		return fmt.Errorf("url cannot be empty")
	}
	
	// Validate URL format
	parsedUrl, err := url.ParseRequestURI(img.Url)
	if err != nil {
		return fmt.Errorf("url is not valid: %w", err)
	}
	
	// Ensure URL has a scheme (http/https)
	if parsedUrl.Scheme == "" {
		return fmt.Errorf("url must include a scheme (http:// or https://)")
	}
	
	// Validate Status is one of the allowed values
	allowedStatuses := map[string]bool{
		"pending":    true,
		"in process": true,
		"done":       true,
		"failed":     true,
	}
	if !allowedStatuses[img.Status] {
		return fmt.Errorf("status must be one of: pending, 'in process', done, failed")
	}
	
	return nil
}

func (img *Image) CanDelete() bool {
	return img.Status != "in process"
}

func (img *Image) MarkAsProcessing() {
	img.Status = "in process"
	img.UpdatedAt = time.Now().UTC()
}

func (img *Image) MarkAsDone(width, height int16, format string) {
	img.Status = "done"
	img.UpdatedAt = time.Now().UTC()
	img.Width = &width
	img.Height = &height
	img.Format = &format
}

func (img *Image) MarkAsFailed() {
	img.Status = "failed"
	img.UpdatedAt = time.Now().UTC()
}

func (img *Image) IsPending() bool {
	return img.Status == "pending"
}