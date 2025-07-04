package gemini

import (
	"context"
	"fmt"
	"os"
	"strings"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func Ask(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("error creating API client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
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