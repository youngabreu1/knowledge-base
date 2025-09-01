package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/youngabreu1/knowledge-base.git/internal/database"
	"github.com/youngabreu1/knowledge-base.git/internal/gemini"
	"github.com/youngabreu1/knowledge-base.git/internal/repository"
	"github.com/youngabreu1/knowledge-base.git/internal/weaviate"
)

func main() {
	log.Println("Starting data ingestion process from database...")
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

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

	err = weaviateService.EnsureClassExists(ctx, "Page")
	if err != nil {
		log.Fatalf("Failed to ensure Weaviate class exists: %v", err)
	}

	pageRepo := repository.NewPageRepository(db)

	pages, err := pageRepo.FindAll()
	if err != nil {
		log.Fatalf("Failed to find all pages from repository: %v", err)
	}
	log.Printf("Found %d pages in the database to process.", len(pages))

	for _, page := range pages {

		log.Printf("Processing page: '%s'", page.Name)

		vector, err := geminiService.CreateEmbedding(ctx, page.Text)
		if err != nil {
			log.Printf("WARNING: could not create embedding for page '%s': %v", page.Name, err)
			continue
		}

		err = weaviateService.AddPage(ctx, "Page", page.Text, page.Name, vector)
		if err != nil {
			log.Printf("WARNING: could not add page '%s' to weaviate: %v", page.Name, err)
		} else {
			log.Printf("Successfully added page '%s' to Weaviate.", page.Name)
		}
	}

	log.Println("Data ingestion process finished.")
}