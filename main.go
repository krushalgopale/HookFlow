package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/krushalgopale/HookFlow/internal/ai"
	"github.com/krushalgopale/HookFlow/internal/broker"
	"github.com/krushalgopale/HookFlow/internal/database"
	"github.com/krushalgopale/HookFlow/internal/middleware"
	"github.com/krushalgopale/HookFlow/internal/routes"
)

func main() {
	// Loading env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database initialization
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres SQL: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	// RabbitMQ initialization
	conn, err := broker.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to RabbitMQ")

	// Queue declaration
	err = broker.DeclareDeliveryQueue(conn)
	if err != nil {
		log.Fatal(err)
	}

	// Consumer initialization
	err = broker.ConsumeDeliveries(conn, db)
	if err != nil {
		log.Fatal(err)
	}

	// AI initialization  
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		log.Fatal("AI_SERVICE_URL is not set")
	}

	aiClient := ai.NewClient(aiURL)

	// Routes
	mux := routes.Routes(db, conn, aiClient)

	// Logger
	logger := middleware.Logger(mux)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: logger,
	}

	// Listen for Ctrl+C and termination signals
	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Starting server
	serverErr := make(chan error, 1)
	go func() {
		log.Println("Starting server on :8080")
		serverErr <- server.ListenAndServe()
	}()

	// Wait for either a shutdown signal or a server error.
	select {
	case <-signalCtx.Done():
		log.Println("Shutdown signal received")

	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server failed: %v", err)
		}
	}

	// Give active HTTP requests up to 10 seconds to finish
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	shutErr := server.Shutdown(shutdownCtx)
	if shutErr != nil {
		log.Printf("Graceful HTTP shutdown failed: %v", err)

		closeErr := server.Close()
		if closeErr != nil {
			log.Printf("HTTP server close failed: %v", closeErr)
		}
	}

	log.Println("HTTP server stopped")
}
