package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/y-maeda1116/template-go-cross/internal/config"
)

func main() {
	cfg := config.Load()
	_ = cfg // Use cfg to prevent unused variable warning

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Println("Application started. Press Ctrl+C to stop.")

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
