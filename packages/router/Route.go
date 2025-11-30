package router

import "github.com/go-chi/chi/v5"

type Routes interface {
	Register(router chi.Router)
}

func SetupRoutes(rgList []Routes, r chi.Router) {
	for _, rg := range rgList {
		rg.Register(r)
	}
}
