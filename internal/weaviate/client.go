package weaviate

import (
    "context"
    "fmt"

    "github.com/weaviate/weaviate-go-client/v4/weaviate"
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