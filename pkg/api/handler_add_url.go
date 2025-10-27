package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/models"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	ImageQueue chan models.ImageProcess
}

// Handler that add url to the database
func (h Handler) HandlerAddUrl(w http.ResponseWriter, r *http.Request) {

	type paramaters struct {
		Url string `json:"url"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err := decoder.Decode(&params)


	if err != nil {
		log.Println("Error decoding request body", err)
		responses.RespondWithError(w, 500, "Error decoding request body")
		return
	}
	
	if len(params.Url) == 0 {
		log.Println("Url cannot be empty")
		responses.RespondWithError(w, 500, "Url cannot be empty")
		return
	}

	parsedUrl, err := url.ParseRequestURI(params.Url)

	if err != nil {
		log.Println("Url not valid", err)
		responses.RespondWithError(w, 500, fmt.Sprintf("Url not valid: %v", err))
		return
	}

	imageProcess := models.ImageProcess{
		ID: 			 	uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		Url: 				parsedUrl.String(),
		Status: 		"pending",
		Width: 			sql.NullInt16{},
		Height: 		sql.NullInt16{},
		Format: 		sql.NullString{},
	}

	err = h.DB.Create(&imageProcess).Error
	if err != nil {
		log.Println("Error creating record", err)
		responses.RespondWithError(w, 500, fmt.Sprintf("Error creating record: %v", err))
		return
	}

	//Add image to queue
	// h.ImageQueue <- imageProcess

	responses.RespondWithJson(w, 200, imageProcess.ToResponse())
}
