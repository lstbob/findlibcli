package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type LLMConfig struct {
	Provider string
	APIKey   string
	Model    string
}

type LLMResolver struct {
	cfg LLMConfig
}

func NewLLM(cfg LLMConfig) *LLMResolver {
	model := cfg.Model
	if model == "" {
		switch cfg.Provider {
		case "gemini":
			model = "gemini-2.0-flash"
		case "openai":
			model = "gpt-4o-mini"
		case "anthropic":
			model = "claude-3-haiku-20240307"
		case "groq":
			model = "llama-3.3-70b-versatile"
		default:
			model = "gemini-2.0-flash"
		}
	}
	return &LLMResolver{cfg: LLMConfig{
		Provider: cfg.Provider,
		APIKey:   cfg.APIKey,
		Model:    model,
	}}
}

type llmLibraryResult struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImportPath  string  `json:"import_path"`
	Confidence  float64 `json:"confidence"`
}

func (r *LLMResolver) Resolve(language, description string) ([]Result, error) {
	switch r.cfg.Provider {
	case "gemini":
		return r.resolveGemini(language, description)
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", r.cfg.Provider)
	}
}

func (r *LLMResolver) resolveGemini(language, description string) ([]Result, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(r.cfg.APIKey))
	if err != nil {
		return nil, fmt.Errorf("creating gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(r.cfg.Model)

	prompt := fmt.Sprintf(`You are a library recommendation engine. Given a programming language and a description of what the user needs, suggest up to 5 libraries or built-in APIs that match.

Language: %s
Need: %s

Return ONLY a JSON array of objects with these fields:
- name: the library or package name
- description: one-sentence summary of what it does
- import_path: the Go import path, npm package name, pip package name, etc (the string the user would use to import/add the dependency)
- confidence: a number from 0.0 to 1.0 indicating how well it matches

Example:
[{"name": "resty", "description": "Simple HTTP client with automatic retries", "import_path": "github.com/go-resty/resty/v2", "confidence": 0.95}]

ONLY output the JSON array, nothing else.`, language, description)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gemini api call: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no response from gemini")
	}

	var textParts []string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			textParts = append(textParts, string(txt))
		}
	}
	text := strings.Join(textParts, "")

	clean := cleanJSON(text)

	var libs []llmLibraryResult
	if err := json.Unmarshal([]byte(clean), &libs); err != nil {
		return nil, fmt.Errorf("parsing gemini response: %w\nraw: %s", err, clean)
	}

	results := make([]Result, 0, len(libs))
	for _, l := range libs {
		results = append(results, Result{
			Library: Library{
				Name:        l.Name,
				Description: l.Description,
				ImportPath:  l.ImportPath,
				Language:    language,
			},
			Confidence: l.Confidence,
		})
	}
	return results, nil
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	return s
}
