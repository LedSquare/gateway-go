package app_server

import (
	"gateway-go/packages/config"
	"log"
	"net/http"
	"time"
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

	return &http.Server{
		Addr:         cfg.GetString("server.address"),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  idleTimeout,
	}
}
