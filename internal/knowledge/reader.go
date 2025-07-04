package knowledge

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LoadFromDirectory(dirPath string) (string, error) {
	var knowledgeBase strings.Builder

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Warning: failed to read file %s. Error: %v", path, err)
				return nil
			}

			knowledgeBase.WriteString(fmt.Sprintf("\n--- Start of Manual: %s ---\n", info.Name()))
			knowledgeBase.Write(content)
			knowledgeBase.WriteString(fmt.Sprintf("\n--- End of Manual: %s ---\n", info.Name()))
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to walk through knowledge directory: %w", err)
	}

	if knowledgeBase.Len() == 0 {
		return "", fmt.Errorf("no .txt files found in directory: %s", dirPath)
	}

	return knowledgeBase.String(), nil
}

func RetrieveContext(prompt, knowledgeBase string) string {
	keywords := strings.Fields(strings.ToLower(prompt))

	paragraphs := strings.Split(knowledgeBase, "\n")

	var relevantContext strings.Builder
	addedParagraphs := make(map[string]bool)

	for _, p := range paragraphs {
		if len(strings.TrimSpace(p)) == 0 || strings.Contains(p, "---") {
			continue
		}

		pLower := strings.ToLower(p)
		
		for _, keyword := range keywords {
			if len(keyword) > 3 && strings.Contains(pLower, keyword) {
				if !addedParagraphs[p] {
					relevantContext.WriteString(p)
					relevantContext.WriteString("\n")
					addedParagraphs[p] = true
				}
			}
		}
	}

	return relevantContext.String()
}