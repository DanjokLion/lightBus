package main

import (
	"context"
	// "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DanjokLion/lightBus/internal/broker"
	myhttp "github.com/DanjokLion/lightBus/internal/http"
	"github.com/DanjokLion/lightBus/pkg/logger"
)

func main() {
	log := logger.New()

	bus := broker.NewInMemoryBroker()

	server := myhttp.NewServer(bus, log)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Starting HTTP server on :8080")
		if err := server.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error: %v", err)
		}

	}()

	<-stop

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Graceful shutdown failed: %v", err)
	} else {
		log.Info("Shutdown complete.")
	}
}