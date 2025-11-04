package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/api"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/config"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/database"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/health"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/middleware"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/worker"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)


func main() {


	conf, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Couldn't load config file with error: %v", err)
	}

	db, err := database.Connect(conf)
	
	if err != nil {
		log.Fatalf("Connection failed with database with error: %v", err)
	}
	
	store := database.NewStore(db)
	repository := store.ImageRepository()
	commander := database.NewImageCommander(store)

	go worker.StartExtract(repository, repository, time.Minute, 5)

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
	apiHandler := api.Handler{
		Querier:   repository,
		Commander: commander,
	}

	v1Router := chi.NewRouter()

	v1Router.Use(middleware.AuthMiddleware(conf.API_KEY))

	apiHandler.Routes(v1Router)

	router.Mount("/v1", v1Router)

	startServer(router, conf.Port)
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