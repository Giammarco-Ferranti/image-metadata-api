package worker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"  // register gif decoder
	_ "image/jpeg" // register jpeg decoder
	_ "image/png"  // register png decoder
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
)

//This worker needs to get from the DB the last N of images.
//Then needs to extract the height/width/format for each image from the url

//start with a startExtract function
func StartExtract(querier domain.ImageQuerier, repository domain.ImageRepository, timeBetweenRequest time.Duration, concurrency int) {
	log.Println("Starting Background Worker")
	ticker := time.NewTicker(timeBetweenRequest)
	ctx := context.Background()

	for ; ; <- ticker.C {
		images, err := querier.FindPendingImages(ctx, concurrency)
	
		if err != nil {
			log.Println("Error retrieving images:", err)
			continue
		}
	
		wg := &sync.WaitGroup{}
		//Process each image
		for _, img := range images {
			wg.Add(1)
			go ProcessImage(img, wg, repository)
		}
		wg.Wait()
	}

}

//Function to process each image
func ProcessImage(img *domain.Image, wg *sync.WaitGroup, repository domain.ImageRepository) {
	defer wg.Done()
	ctx := context.Background()

	// Mark as processing using domain method
	img.MarkAsProcessing()
	err := repository.Save(ctx, img)

	if err != nil {
		errorParse(img, "Error updating the image status for image id", err, repository)
		return
	}

	imageData, err := ProcessRequest(img.Url)
	if err != nil {
		errorParse(img, "Error request url for image id", err, repository)
		return 
	}

	width, height, format, err := DecodeImage(imageData)
	if err != nil {
		errorParse(img, "Couldn't decode image id", err, repository)
		return
	}

	// Mark as done using domain method
	img.MarkAsDone(width, height, format)
	err = repository.Save(ctx, img)

	if err != nil {
		errorParse(img, "Couldn't Update image with status done for image id", err, repository)
		return
	}

	log.Println("Successfully processed image id: ", img.ID)
	log.Println("============================")
	
}

func errorParse(img *domain.Image, msg string, err error, repository domain.ImageRepository) {
	log.Printf("%v: %v, with error: %v", msg, img.ID, err)
	ctx := context.Background()
	img.MarkAsFailed()
	saveErr := repository.Save(ctx, img)
	if saveErr != nil {
		log.Printf("Error updating failed image for img id: %v, with error: %v", img.ID, saveErr)
	}
}

func ProcessRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func DecodeImage(imageData []byte) (int16, int16, string, error) {
	imgValue, format, err := image.Decode(bytes.NewReader(imageData))

	if err != nil {
		return 0, 0, "", err
	}

	bounds := imgValue.Bounds()
	width := int16(bounds.Dx()) //width in pixels
	height := int16(bounds.Dy()) // height in pixels

	return width, height, format, nil
}