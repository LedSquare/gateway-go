package router

import "github.com/go-chi/chi/v5"

type Routes interface {
	Register(router chi.Router)
}
