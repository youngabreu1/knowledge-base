package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/youngabreu1/knowledge-base.git/internal/gemini"
)

type ApiHandler struct {
	geminiService *gemini.GeminiService
}

func NewApiHandler(geminiSvc *gemini.GeminiService) *ApiHandler {
	return &ApiHandler{
		geminiService: geminiSvc,
	}
}

func (h *ApiHandler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		http.Error(w, "Error: 'prompt' parameter is required.", http.StatusBadRequest)
		return
	}
	log.Printf("Prompt received: \"%s\"", prompt)

	context := "Any context was found."

	finalPrompt := fmt.Sprintf(`
		You will be the AI ​​agent for Playlist Software Solutions. You will be the first contact to try to solve a problem sent by the customer. Another AI Agent will send you the costumer question and you have to try to solve it. This message has a context in which we searched our internal documentation to help you answer the question. Based on this document, reply the costumer question, if you're not able to reply the question, say that isn't in your documentation. Do not provide information there's not in your context. Answer in PT-BR

		Context:
		%s

		Question:
		%s
	`, context, prompt)

	response, err := h.geminiService.Ask(r.Context(), finalPrompt)
	if err != nil {
		log.Printf("Error getting response from Gemini: %v", err)
		http.Error(w, "Failed to get a response from the AI agent.", http.StatusInternalServerError)
		return
	}

	log.Println("Context-aware response sent successfully!")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprint(w, response)
}

