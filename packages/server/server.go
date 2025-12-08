package app_server

import (
	"gateway-go/internal/routes"
	"gateway-go/packages/config"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func Make() *http.Server {
	cfg := config.Load()

	timeout, err := time.ParseDuration(cfg.GetString("server.timeout"))
	if err != nil {
		log.Fatalf("wrong format server.timeout: %v", err)
	}

	idleTimeout, err := time.ParseDuration(cfg.GetString("server.idle_timeout"))
	if err != nil {
		log.Fatalf("wrong format server.idle_timeout: %v", err)
	}

	// pp.Print(cfg.GetStringSlice("proxy.allowed_headers"))
	r := chi.NewRouter()
	routes.Register(r)

	return &http.Server{
		Addr:         cfg.GetString("server.address"),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  idleTimeout,
		// Handler:      r,
	}
}
