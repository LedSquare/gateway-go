package routes

import (
	"gateway-go/internal/routes/api"

	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router) {
	api.V1(r)
}
