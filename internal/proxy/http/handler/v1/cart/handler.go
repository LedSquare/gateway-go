package cart

import (
	"gateway-go/internal/proxy/http/handler/v1/proxy"
	"gateway-go/packages/config"
	"net/http"
)

type CartHandler struct {
	handler *proxy.Handler
}

func NewCartHandler() *CartHandler {
	cfg := config.Load()
	return &CartHandler{
		proxy.NewController(
			cfg.GetString("proxy.cart.url"),
			cfg.GetStringSlice("proxy.allowed_headers"),
		),
	}
}

func (h *CartHandler) Send(w http.ResponseWriter, r *http.Request) {
	if err := h.handler.Send(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}
