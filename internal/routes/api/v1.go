package api

import (
	handlers "gateway-go/internal/proxy/http/handlers/v1"

	"github.com/go-chi/chi/v5"
)

func V1(r chi.Router) {

	r.Route("/v1", func(r chi.Router) {
		cartHandler := handlers.NewCartHandler()
		r.Route("/cart", func(r chi.Router) {
			r.MethodNotAllowed(nil)
			r.HandleFunc("/*", cartHandler.Send)
		})
	})
}
