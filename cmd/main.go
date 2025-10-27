package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/api"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/health"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/middleware"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/worker"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)


func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Printf("Couldn't find the port")
	}

	// Build database connection string from environment variables
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatal("Missing required database environment variables")
	}

	databaseUrl := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"

	db, err := gorm.Open(postgres.Open(databaseUrl))
	if err != nil {
		log.Printf("Connection failed with database with error: %v", err)
	}

	db.AutoMigrate(&models.ImageProcess{})

	go worker.StartExtract(db, time.Minute)

	log.Printf("Connection database started")

	router := chi.NewRouter()


	//Configure CORS
	router.Use(cors.Handler(cors.Options{
    // AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
    AllowedOrigins:   []string{"https://*", "http://*"},
    // AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: false,
    MaxAge:           300, // Maximum value not ignored by any of major browsers
  }))

	//Health endpoint
	router.Get("/healthz", health.HandlerHealth)

	//Initiate api handler
	apiHandler := api.Handler{DB: db}

	v1Router := chi.NewRouter()

	v1Router.Use(middleware.AuthMiddleware)

	v1Router.Get("/images", apiHandler.HandlerGetImages)

	v1Router.Post("/image", apiHandler.HandlerAddUrl)
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