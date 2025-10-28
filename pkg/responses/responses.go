package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
)

//For list responses
type ListResponse struct {
	Data []models.ImageProcessResponse `json:"data"`
}

//For single item responses
type ItemResponse struct {
	Data models.ImageProcessResponse 		`json:"data"`
}
//For responses with metadata
type PaginatedResponse struct {
	Data []models.ImageProcessResponse `json:"data"`
	Meta PaginationMeta                `json:"meta"`
}

type PaginationMeta struct {
	Total  int 	`json:"total"`
	Offset int 	`json:"offset"`
	ResultsCount  int	`json:"results_count"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error:", msg)
	}

	type ErrorMsg struct {
		Error string `json:"error"`
	}

	RespondWithJson(w, code, ErrorMsg{Error: msg})
}


func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

