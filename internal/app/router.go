package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kfreiman/imageservice/pkg/handler"
)

// NewRouter returns handler for Server's HTTP implementation
func NewRouter(handler handler.HTTPHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/image", func(r chi.Router) {
		r.Get("/modify", handler.ModifyEndpoint())
		r.Post("/modify", handler.ModifyEndpoint())
	})

	// other app-level endpoints. It can be prometeus, pprof, healthchecks etc..

	return r
}
