package gemini

import (
	"context"
	"fmt"
	"os"
    "strings"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService struct {
    client          *genai.Client
    generativeModel *genai.GenerativeModel
    embeddingModel  *genai.EmbeddingModel
}

func NewService(ctx context.Context) (*GeminiService, error) {  
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
    }

    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, fmt.Errorf("error creating API client: %w", err)
    }

    generativeModel := client.GenerativeModel("gemini-1.5-flash")
    embeddingModel := client.EmbeddingModel("text-embedding-004")

    return &GeminiService{
        client:          client,
        generativeModel: generativeModel,
        embeddingModel:  embeddingModel,
    }, nil

}

func (s *GeminiService) Close() {
    s.client.Close()
}

func (s *GeminiService) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {

    res, err := s.embeddingModel.EmbedContent(ctx, genai.Text(text))
    if err != nil {
        return nil, fmt.Errorf("error creating embedding: %w", err)
    }
    if res.Embedding == nil || res.Embedding.Values == nil {
        return nil, fmt.Errorf("received an empty embedding from the API")
    }
    return res.Embedding.Values, nil
}

func (s *GeminiService) Ask(ctx context.Context, prompt string) (string, error) {
    resp, err := s.generativeModel.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        return "", fmt.Errorf("error generating content: %w", err)
    }
    return parseResponse(resp), nil
}

func parseResponse(resp *genai.GenerateContentResponse) string {
	var responseText strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				responseText.WriteString(fmt.Sprintf("%s", part))
			}
		}
	}
	return responseText.String()
}