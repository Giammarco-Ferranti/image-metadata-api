package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"
	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (h Handler) HandleDeleteImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	imageIdString := chi.URLParam(r, "id")
	imageId, err := uuid.Parse(imageIdString)

	if err != nil {
		slog.Error("Couldn't parse uuid", "error", err)
		responses.RespondWithError(w, 400, fmt.Sprintf("Couldn't parse id: %v", err))
		return
	}


	image, err := h.Querier.FindById(ctx, imageId)

	if err != nil {
		slog.Error("Couldn't get image", "error", err)
		responses.RespondWithError(w, 500, fmt.Sprintf("Couldn't get image: %v", err))
		return
	}

	if image == nil {
		responses.RespondWithError(w, 404, "Image not found")
		return
	}

	if !image.CanDelete() {
		slog.Error("Cannot delete image with status 'in process'")
		responses.RespondWithError(w, 400, "Cannot delete image with status 'in process'")
		return
	}


	repo, ok := h.Querier.(domain.ImageRepository)
	if !ok {
		slog.Error("Querier is not an ImageRepository")
		responses.RespondWithError(w, 500, "Internal server error")
		return
	}

	err = repo.Delete(ctx, imageId)
	if err != nil {
		slog.Error("Error deleting image", "error", err)
		responses.RespondWithError(w, 500, fmt.Sprintf("Error deleting image: %v", err))
		return
	}

	responses.RespondWithJson(w, 200, "Successfully deleted image")

}