package api

import "github.com/go-chi/chi/v5"

func V1(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		r.Get("/example")
	})
}
