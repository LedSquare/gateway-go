package main

import (
	"context"
	appServer "gateway-go/packages/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	srv := appServer.Make()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server is down: %s", err.Error())
		}
	}()

	log.Print("Server start!")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	log.Print("Shutting down server")

	// Даём 10 секунд на graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Print("Server stopped")
}
