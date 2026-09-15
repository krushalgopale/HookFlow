package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
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
		log.Fatal(err)
	}
	defer db.Close()
		
	log.Println("Connected to PostgreSQL")
	
	// Routes
	mux := routes.Routes(db)
 
	logger := middleware.Logger(mux)

	log.Println("Starting server on :8080")

	// Starting server
	http.ListenAndServe("0.0.0.0:8080", logger)
}
