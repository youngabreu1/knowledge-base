package weaviate

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/youngabreu1/knowledge-base.git/internal/gemini" 
	"strings"
)

type WeaviateService struct {
	Client *weaviate.Client
}

func NewService(host, scheme string) (*WeaviateService, error) {
	cfg := weaviate.Config{
		Host:   host,
		Scheme: scheme,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create weaviate client: %w", err)
	}

	ready, err := client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("weaviate not ready: %w", err)
	}
	if !ready {
		return nil, fmt.Errorf("weaviate is not ready to accept connections")
	}

	return &WeaviateService{
		Client: client,
	}, nil
}

func (s *WeaviateService) AddPage(ctx context.Context, className string, content string, source string, vector []float32) error {
	properties := map[string]interface{}{
		"content": content,
		"source":  source,
	}

	_, err := s.Client.Data().Creator().
		WithClassName(className).
		WithID(uuid.New().String()).
		WithProperties(properties).
		WithVector(vector).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to add page to weaviate: %w", err)
	}

	return nil
}

func (s *WeaviateService) EnsureClassExists(ctx context.Context, className string) error {
	_, err := s.Client.Schema().ClassGetter().WithClassName(className).Do(ctx)
	if err == nil {
		log.Printf("Class '%s' already exists. Skipping creation.", className)
		return nil
	}

	log.Printf("Class '%s' not found. Creating it now.", className)

	classObj := &models.Class{
		Class:       className,
		Description: "Stores pages of documentation from the knowledge base.",
		Vectorizer:  "none",
		Properties: []*models.Property{
			{
				Name:        "content",
				DataType:    []string{"text"},
				Description: "The actual text content of the page chunk.",
			},
			{
				Name:        "source",
				DataType:    []string{"text"},
				Description: "The original source or title of the document.",
			},
		},
	}

	err = s.Client.Schema().ClassCreator().WithClass(classObj).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class '%s': %w", className, err)
	}

	log.Printf("Successfully created class '%s'.", className)
	return nil
}

func (s *WeaviateService) SearchPages(ctx context.Context, className string, prompt string, geminiSvc *gemini.GeminiService) (string, error) {
    vector, err := geminiSvc.CreateEmbedding(ctx, prompt)
    if err != nil {
        return "", fmt.Errorf("failed to create embedding for prompt: %w", err)
    }

    content := graphql.Field{Name: "content"}
    source := graphql.Field{Name: "source"}

    nearVector := s.Client.GraphQL().NearVectorArgBuilder().
        WithVector(vector)

    response, err := s.Client.GraphQL().Get().
        WithClassName(className).
        WithFields(content, source).
        WithNearVector(nearVector).
        WithLimit(39).
        Do(ctx)

    if err != nil {
        return "", fmt.Errorf("failed to perform search in weaviate: %w", err)
    }

    data := response.Data["Get"].(map[string]interface{})
    pages := data["Page"].([]interface{})

    var contextBuilder strings.Builder
    for _, item := range pages {
        pageMap := item.(map[string]interface{})
        if contentStr, ok := pageMap["content"].(string); ok {
            contextBuilder.WriteString(contentStr)
            contextBuilder.WriteString("\n---\n")
        }
    }
    
    return contextBuilder.String(), nil
}

// ---------------------- DEBUG ONLY ---------------------------
// func (s *WeaviateService) GetSchema(ctx context.Context) error {
// 	schema, err := s.Client.Schema().Getter().Do(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to get schema: %v", err)
// 	}
// 	fmt.Println("Weaviate Classes: ")
// 	for _, class := range schema.Classes {
// 		fmt.Printf("- %s\n", class.Class)
// 	}
// 	return err
// }

// func (s *WeaviateService) Iterator(ctx context.Context) error {
//     iterator := s.Client.Data().ObjectsGetter().WithClassName("Page").WithLimit(100)
//     for {
//         objects, err := iterator.Do(ctx)
//         if err != nil {
//             log.Fatalf("Failed iterating data: %v", err)
//         }
//         for _, obj := range objects {
//             fmt.Printf("obj: %v", obj)
//         }
//     }
// }