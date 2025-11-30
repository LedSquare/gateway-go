package router

import "github.com/go-chi/chi/v5"

type Route interface {
	Register(router chi.Router)
}
