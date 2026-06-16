package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/lstbob/findlibcli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "View or set configuration",
	Long: `View or set configuration values.

Supported keys:
  gemini       Set Gemini API key (free: aistudio.google.com/apikey)
  openai       Set OpenAI API key
  anthropic    Set Anthropic API key
  groq         Set Groq API key
  context7     Set Context7 API key (free: context7.com/dashboard)
  llm          Set LLM provider (gemini, openai, anthropic, groq)

Examples:
  findlib config               # show current config
  findlib config gemini sk-xxx # set Gemini API key
  findlib config llm groq      # switch LLM provider to Groq
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		switch len(args) {
		case 0:
			printConfig(cfg)
		case 2:
			setConfig(cfg, args[0], args[1])
		default:
			fmt.Fprintln(os.Stderr, "Usage: findlib config [key] [value]")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func printConfig(cfg *config.Config) {
	fmt.Println("Current configuration:")
	fmt.Printf("  LLM Provider:  %s\n", cfg.LLMProvider)
	masked := func(s string) string {
		if s == "" {
			return "(not set)"
		}
		if len(s) > 8 {
			return s[:4] + "…" + s[len(s)-4:]
		}
		return "****"
	}
	fmt.Printf("  Gemini API:    %s\n", masked(cfg.GeminiAPIKey))
	fmt.Printf("  OpenAI API:    %s\n", masked(cfg.OpenAIAPIKey))
	fmt.Printf("  Anthropic API: %s\n", masked(cfg.AnthropicAPIKey))
	fmt.Printf("  Groq API:      %s\n", masked(cfg.GroqAPIKey))
	fmt.Printf("  Context7 API:  %s\n", masked(cfg.Context7APIKey))
	fmt.Println()
	fmt.Println("Config file: ~/.config/findlib/config.json")
}

func setConfig(cfg *config.Config, key, value string) {
	switch strings.ToLower(key) {
	case "gemini":
		cfg.GeminiAPIKey = value
	case "openai":
		cfg.OpenAIAPIKey = value
	case "anthropic":
		cfg.AnthropicAPIKey = value
	case "groq":
		cfg.GroqAPIKey = value
	case "context7":
		cfg.Context7APIKey = value
	case "llm":
		cfg.LLMProvider = value
	default:
		fmt.Fprintf(os.Stderr, "Unknown config key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Valid keys: gemini, openai, anthropic, groq, context7, llm\n")
		os.Exit(1)
	}

	if err := config.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ %s updated\n", key)
}
