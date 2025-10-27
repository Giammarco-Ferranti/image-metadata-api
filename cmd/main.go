package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/api"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/health"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/middleware"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/worker"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)


func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Printf("Couldn't find the port")
	}

	databaseUrl := os.Getenv("DB_URL")
	if databaseUrl == "" {
		log.Printf("Couldn't find the databaseUrl")
	}

	db, err := gorm.Open(postgres.Open(databaseUrl))
	if err != nil {
		log.Printf("Connection failed with database with error: %v", err)
	}

	db.AutoMigrate(&models.ImageProcess{})

	log.Printf("Connection database started")

	router := chi.NewRouter()

	imageQueue := make(chan models.ImageProcess, 10)

	go worker.StartProcessImage(db, imageQueue)

	//Configure CORS

	//Health endpoint
	router.Get("/healthz", health.HandlerHealth)

	//Initiate api handler
	apiHandler := api.Handler{DB: db, ImageQueue: imageQueue}

	v1Router := chi.NewRouter()

	v1Router.Use(middleware.AuthMiddleware)
	v1Router.Post("/url", apiHandler.HandlerAddUrl)

	//Retrieve all images
	v1Router.Get("/images", apiHandler.HandlerGetImages)

	//Retrieve single image
	v1Router.Get("/image/{id}", apiHandler.HandlerGetImage)

	router.Mount("/v1", v1Router)

	startServer(router, portString)
}

func startServer(router chi.Router, portString string) {
	server := &http.Server{
		Addr: ":" + portString,
		Handler: router,
	}

	log.Printf("Server started with port: %v", portString)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}