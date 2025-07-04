package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/youngabreu1/knowledge-base.git/internal/gemini"
	"github.com/youngabreu1/knowledge-base.git/internal/knowledge"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// query?prompt=something
	knowledge, err := knowledge.LoadFromDirectory("./guides")
	if err !=nil {
		log.Printf("Error loading knowlodge base: ", err)
		http.Error(w,"Could not load knowledge base.", http.StatusInternalServerError)
		return
	} 
	
	log.Println("Knowledge base loaded succesfully...")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, knowledge)

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