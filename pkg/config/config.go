package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
Port string
DB_URL string
API_KEY string
}

func LoadConfig() (*Config, error) {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		return nil, fmt.Errorf("couldn't find the port")
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || dbPort == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	databaseUrl := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing required api key")
	}

	return &Config{
		Port: portString,
		DB_URL: databaseUrl,
		API_KEY: apiKey,
	}, nil
}