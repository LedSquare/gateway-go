package api

import (
	"gateway-go/internal/proxy/http/handler/v1/cart"

	"github.com/go-chi/chi/v5"
)

func V1(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		ch := cart.NewCartHandler()
		r.Get("/example", ch.Send())
	})
}
