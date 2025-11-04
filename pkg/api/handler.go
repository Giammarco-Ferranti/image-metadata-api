package api

import "github.com/Giammarco-Ferranti/image-metadata-api/pkg/domain"

type Handler struct {
	Querier   domain.ImageQuerier
	Commander domain.ImageCommander
}