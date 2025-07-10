package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/youngabreu1/knowledge-base.git/internal/database"
	"github.com/youngabreu1/knowledge-base.git/cmd/api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: Could not load .env file.")
	}

	http.HandleFunc("/query", api.ApiHandler)

	port := os.Getenv("PORT")
	log.Printf("Server started. Listening on port %s", port)
	log.Printf("Endpoint available at: http://localhost:%s/query", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Fatal error starting server: %v", err)
	}
}