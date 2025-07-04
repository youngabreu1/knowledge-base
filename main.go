package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/youngabreu1/knowledge-base.git/internal/gemini"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// query?prompt=something
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		http.Error(w, "Error: 'prompt' parameter is required.", http.StatusBadRequest)
		return
	}
	log.Printf("Prompt received: \"%s\"", prompt)

	response, err := gemini.Ask(prompt)
	if err != nil {
		log.Printf("Error calling gemini package: %v", err)
		http.Error(w, "Error processing request with Gemini.", http.StatusInternalServerError)
		return
	}

	log.Println("Response sent successfully!")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: Could not load .env file.")
	}

	http.HandleFunc("/query", apiHandler)

	port := os.Getenv("PORT")
	log.Printf("Server started. Listening on port %s", port)
	log.Printf("Endpoint available at: http://localhost:%s/query", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Fatal error starting server: %v", err)
	}
}