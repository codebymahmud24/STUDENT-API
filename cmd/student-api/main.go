package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codebymahmud24/student-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// Database Setup

	// Setup Router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// Start Server
	server := &http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: router,
	}

	// Channel to listen for interrupt signal to terminate gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server is running on http://localhost%s\n", cfg.HttpServer.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Block until we receive our signal
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for graceful shutdown (30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
