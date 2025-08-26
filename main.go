package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/youngabreu1/knowledge-base.git/cmd/api"
	"github.com/youngabreu1/knowledge-base.git/internal/database"
	"github.com/youngabreu1/knowledge-base.git/internal/gemini"
	"github.com/youngabreu1/knowledge-base.git/internal/repository"
	"github.com/youngabreu1/knowledge-base.git/internal/weaviate"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: Could not load .env file.")
	}
	port := os.Getenv("PORT")

	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed establishing connection to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection succeed.")

	geminiService, err := gemini.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v", err)	
	}
	defer geminiService.Close()

    weaviateHost := os.Getenv("WEAVIATE_HOST")
    weaviateScheme := os.Getenv("WEAVIATE_SCHEME")
    
    weaviateService, err := weaviate.NewService(weaviateHost, weaviateScheme)
    if err != nil {
        log.Fatalf("Failed to create Weaviate service: %v", err)
    }
    log.Println("Weaviate connection succeed: ", weaviateService)

	 err = weaviateService.EnsureClassExists(ctx, "Page")
    if err != nil {
        log.Fatalf("Failed to ensure Weaviate class exists: %v", err)
    }

	pageRepo := repository.NewPageRepository(db)
	pageRepo.FindAll()
    log.Println("Page repo created successfully.", pageRepo)


	apiHandler := api.NewApiHandler(geminiService)
	http.HandleFunc("/query", apiHandler.HandleQuery)
	log.Printf("Endpoint available at: http://localhost:%s/query", port)
	err = http.ListenAndServe(":"+port, nil)
	if (err != nil) {
		log.Fatalf("Fatal error starting server: %v", err)
	}
	
}