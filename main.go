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
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		http.Error(w, "Error: 'prompt' parameter is required.", http.StatusBadRequest)
		return
	}
	log.Printf("Prompt received: \"%s\"", prompt)

	knowledgeBase, err := knowledge.LoadFromDirectory("./guides")
	if err != nil {
		log.Printf("Error loading knowledge base: %v", err)
		http.Error(w, "Could not load knowledge base.", http.StatusInternalServerError)
		return
	}

	context := knowledge.RetrieveContext(prompt, knowledgeBase)
	if context == "" {
		context = "No relevant context found."
	}

	finalPrompt := fmt.Sprintf(`
		You will be the AI ​​agent for XXX . You will be the first contact to try to solve a problem sent by the customer. This message has a context in which we searched our internal documentation to help you answer the question. Based on this document, reply the costumer question, if you don't know the answer, say that isn't in your documentation. Answer in PT-BR

		Context:
		%s

		Question:
		%s
	`, context, prompt)


	response, err := gemini.Ask(finalPrompt)
	if err != nil {
		log.Printf("Error calling gemini package: %v", err)
		http.Error(w, "Error processing request with Gemini.", http.StatusInternalServerError)
		return
	}

	log.Println("Context-aware response sent successfully!")
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