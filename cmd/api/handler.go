package api

import (
	"fmt"
	"log"
	"net/http"
	"github.com/youngabreu1/knowledge-base.git/internal/weaviate"
	"github.com/youngabreu1/knowledge-base.git/internal/gemini"
)

type ApiHandler struct {
	geminiService *gemini.GeminiService
	weaviateService *weaviate.WeaviateService
}

func NewApiHandler(geminiSvc *gemini.GeminiService, weaviateSvc *weaviate.WeaviateService) *ApiHandler {
	return &ApiHandler{
		geminiService: geminiSvc,
		weaviateService: weaviateSvc,
	}
}

func (h *ApiHandler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		http.Error(w, "Error: 'prompt' parameter is required.", http.StatusBadRequest)
		return
	}
	log.Printf("Prompt received: \"%s\"", prompt)

    log.Println("Searching for context in Weaviate...")
    context, err := h.weaviateService.SearchPages(r.Context(), "Page", prompt, h.geminiService)
    if err != nil {
        log.Printf("Error searching for context: %v", err)
        http.Error(w, "Failed to search for context.", http.StatusInternalServerError)
        return
    }

    if context == "" {
        context = "Nenhuma informação relevante foi encontrada em nossa base de conhecimento."
    }
    log.Println("Context found.")

	finalPrompt := fmt.Sprintf(
		"You will be the AI agent for Playlist Software Solutions. "+
			"You will be the first contact to try to solve a problem sent by the customer. "+
			"Another AI Agent will send you the customer question, and you have to try to solve it. "+
			"This message has a context in which we searched our internal documentation to help you answer the question. "+
			"Based on this document, reply to the customer question. If you're not able to reply to the question, say that it isn't in your documentation. "+
			"Do not provide information that's not in your context. Answer in PT-BR.\n\n"+
			"Context:\n%s\n\n"+
			"Question:\n%s",
		context, prompt,
	)

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

