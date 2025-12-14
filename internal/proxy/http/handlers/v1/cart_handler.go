package v1_handlers

import (
	"gateway-go/internal/proxy/http/handlers"
	"gateway-go/packages/config"
	"net/http"
)

type CartHandler struct {
	handler *handlers.Handler
}

func NewCartHandler() *CartHandler {
	cfg := config.Load()
	return &CartHandler{
		handlers.NewHandler(
			cfg.GetString("proxy.services.cart.url"),
			cfg.GetStringSlice("proxy.allowed_headers"),
			"cart_handler",
			"v1",
		),
	}
}

func (h *CartHandler) Send(w http.ResponseWriter, r *http.Request) {
	h.handler.Send(w, r)
}
