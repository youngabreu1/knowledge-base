package api

import (
	"fmt"
	"log"
	"net/http"
	"github.com/youngabreu1/knowledge-base.git/internal/gemini"
	"github.com/youngabreu1/knowledge-base.git/internal/knowledge"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) (){
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
		You will be the AI ​​agent for Playlist Software Solutions. You will be the first contact to try to solve a problem sent by the customer. Another AI Agent will send you the costumer question and you have to try to solve it. This message has a context in which we searched our internal documentation to help you answer the question. Based on this document, reply the costumer question, if you're not able to reply the question, say that isn't in your documentation. Answer in PT-BR

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