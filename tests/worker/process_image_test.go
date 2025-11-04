package worker

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/database"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
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
	db.AutoMigrate(&database.ImageModel{})

	store := database.NewStore(db)
	repository := store.ImageRepository()

	testserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(testPng)
	}))

	defer testserver.Close()

	now := time.Now().UTC()
	testImg := &domain.Image{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Url:       testserver.URL,
		Status:    "pending",
		Width:     nil,
		Height:    nil,
		Format:    nil,
	}

	ctx := context.Background()
	err = repository.Create(ctx, testImg)
	if err != nil {
		t.Errorf("Error creating test image: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	worker.ProcessImage(testImg, &wg, repository)
	wg.Wait()

	//Assert the results
	result, err := repository.FindById(ctx, testImg.ID)
	if err != nil {
		t.Errorf("Error finding result: %v", err)
	}

	if result.Status != "done" {
		t.Errorf("Expected status 'done', got %v", result.Status)
	}

}

func TestProcessRequest(t *testing.T) {

	//Create a test http server that returns a fake image
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Write(testPng)
	}))
	defer testServer.Close()

	data, err := worker.ProcessRequest(testServer.URL)

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

