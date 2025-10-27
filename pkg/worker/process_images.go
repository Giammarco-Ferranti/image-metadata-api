package worker

import (
	"log"
	"net/http"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"gorm.io/gorm"
)


func StartProcessImage(db *gorm.DB, imageQueue chan models.ImageProcess) {
	log.Println("Background worker started - checking for pending images every 5 seconds...")
	for {
		log.Println("Checking for pending images...")
		ProcessPendingImages(db, imageQueue)
	}
}


func ProcessPendingImages(db *gorm.DB, queue chan models.ImageProcess) {

	if len(queue) == 0 {
		log.Println("No pending images found")
		return
	}

	log.Printf("Found %d pending image(s) to process", len(queue))

	for image := range queue {
		log.Println("==========================")
		db.Where("id = ?", image.ID).Updates(models.ImageProcess{
			UpdatedAt: time.Now().UTC(),
			Status: "in process",
		})
		resp, err := http.Get(image.Url)
		if err != nil {
			log.Println("Error getting image url", err)
			continue
		}
		log.Println(resp)
	}
}