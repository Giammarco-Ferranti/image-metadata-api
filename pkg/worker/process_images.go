package worker

import (
	"bytes"
	"database/sql"
	"image"
	_ "image/gif"  // register gif decoder
	_ "image/jpeg" // register jpeg decoder
	_ "image/png"  // register png decoder
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"gorm.io/gorm"
)

//This worker needs to get from the DB the last N of images.
//Then needs to extract the height/width/format for each image from the url

//start with a startExtract function
func StartExtract(DB *gorm.DB, timeBetweenRequest time.Duration) {
	log.Println("Starting Background Worker")
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <- ticker.C {
		var images []models.ImageProcess
		err := DB.Where("status = ?", "pending").Order("created_at asc").Limit(5).Find(&images).Error
	
		if err != nil {
			log.Println("Error retrieving images")
			continue
		}
	
		wg := &sync.WaitGroup{}
		//Process each image
		for _, image := range images {
			wg.Add(1)
			go processImage(image, wg, DB)
		}
		wg.Wait()
	}

}

//Function to process each image
func processImage(img models.ImageProcess, wg *sync.WaitGroup, DB *gorm.DB) {
	defer wg.Done()
	err := DB.Where("id = ?", img.ID).Updates(models.ImageProcess{
		UpdatedAt: time.Now().UTC(),
		Status: "in process",
	}).Error

	if err != nil {
		errorParse(img, "Error updating the image status for image id", err, DB)
		return
	}

	resp, err := http.Get(img.Url)

	if err != nil {
		errorParse(img, "Error request url for image id", err, DB)
		return
	}

	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)

	if err != nil {
		errorParse(img, "Couldn't extract image data for image id", err, DB)
		return
	}

	imgValue, format, err := image.Decode(bytes.NewReader(imageData))

	if err != nil {
		errorParse(img, "Couldn't decode image id", err, DB)
		return
	}

	bounds := imgValue.Bounds()
	width := int16(bounds.Dx()) //width in pixels
	height := int16(bounds.Dy()) // height in pixels

	DB.Model(&img).Updates(models.ImageProcess{
		UpdatedAt: time.Now().UTC(),
		Width: sql.NullInt16{Int16: width, Valid: true},
		Height: sql.NullInt16{Int16: height, Valid: true},
		Format: sql.NullString{String: format, Valid: true},
		Status: "done",
	})

	log.Println("Successfully processed image id: ", img.ID)
	log.Println("============================")
	
}

func errorParse(img models.ImageProcess, msg string, err error, DB *gorm.DB) {
	log.Printf("%v: %v, with error: %v", msg, img.ID, err)
	DB.Model(&img).Updates(models.ImageProcess{
		UpdatedAt: time.Now().UTC(),
		Status: "failed",
	})
}