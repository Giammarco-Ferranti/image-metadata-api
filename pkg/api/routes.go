package api

import "github.com/go-chi/chi"

//This routes to different endpoints
func (h Handler) Routes(v1Router *chi.Mux) {
	v1Router.Get("/images", h.HandlerGetImages)

	v1Router.Post("/image", h.HandlerAddImage)
	v1Router.Get("/image/{id}", h.HandlerGetImage)
	v1Router.Delete("/image/{id}", h.HandleDeleteImage)
}