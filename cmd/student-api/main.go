package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codebymahmud24/student-api/internal/config"
	"github.com/codebymahmud24/student-api/internal/http/handler/student"
	"github.com/codebymahmud24/student-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// Database Setup
	storage, err := sqlite.InitialiazeDB(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.ENV), slog.String("version", "1.0.0"))

	// Setup Router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/v1/students", student.CreateStudentHandler(storage))
	router.HandleFunc("GET /api/v1/students/{id}", student.GetStudentById(storage))
	router.HandleFunc("GET /api/v1/students", student.GetAllStudents(storage))
	router.HandleFunc("PUT /api/v1/students/{id}", student.UpdateStudentById(storage))

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
