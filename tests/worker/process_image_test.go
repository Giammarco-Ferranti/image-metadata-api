package worker

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/worker"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed 1x1.png
var testPng []byte

func TestProcessImage(t *testing.T) {

	//Create db connection
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		t.Errorf("Error db connection, %v", err)
	}

	//Migrate schema
	db.AutoMigrate(&models.ImageProcess{})

	testserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(testPng)
	}))

	defer testserver.Close()

	testImg := models.ImageProcess{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url: testserver.URL,
		Status: "pending",
	}

	db.Create(&testImg)

	var wg sync.WaitGroup
	wg.Add(1)
	worker.ProcessImage(testImg, &wg, db)
	wg.Wait()

	//Assert the results
	var result models.ImageProcess
	db.First(&result, testImg.ID)

	if result.Status != "done" {
		t.Errorf("Expected status 'done', got %v", result)
	}

}

func TestProcessRequest(t *testing.T) {

	//Create a test http server that returns a fake image
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Write(testPng)
	}))
	defer testServer.Close()

	img := models.ImageProcess{
		Url: testServer.URL,
	}

	data, err := worker.ProcessRequest(img)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected image data, got empty")
	}
}

func TestDecodeImage(t *testing.T) {
	
	//decode an image with asserted bytes
	height, width, format, _ := worker.DecodeImage(testPng)

	if height != 1 && width != 1 && format != "png" {
		t.Errorf("Expected 1px height and image and png format, got: %v, %v, %v", height, width, format)
	}
}

