package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/lstbob/findlibcli/internal/config"
	"github.com/lstbob/findlibcli/internal/docs"
	"github.com/lstbob/findlibcli/internal/output"
	"github.com/lstbob/findlibcli/internal/resolver"
	"github.com/lstbob/findlibcli/internal/tui"
)

var (
	jsonFlag     bool
	importFlag   bool
	noCopyFlag   bool
)

var rootCmd = &cobra.Command{
	Use:   "findlib [language] [description...]",
	Short: "Find libraries by natural language description",
	Long: `findlib searches for libraries and APIs matching a description in a given language.

Examples:
  findlib go "http client with retries"
  findlib python "json parsing"
  findlib js "state management"
  findlib config gemini "your-api-key"
`,
	Args: cobra.MinimumNArgs(2),
	Run:  run,
}

var searchCmd = &cobra.Command{
	Use:   "search <language> <description...>",
	Short: "Search for libraries matching a description",
	Args:  cobra.MinimumNArgs(2),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	language := args[0]
	description := strings.Join(args[1:], " ")

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	results, err := resolve(cfg, language, description)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("No libraries found matching that description.")
		return
	}

	var lib *resolver.Result
	if jsonFlag || importFlag {
		lib = &results[0]
	} else {
		var pickErr error
		lib, pickErr = tui.PickLibrary(results)
		if pickErr != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", pickErr)
			os.Exit(1)
		}
	}

	docSvc := docs.NewClient(cfg.Context7APIKey)
	docResult, _ := docSvc.GetDocs(lib.Name, description)

	var codeSnippets []string
	var infoText string
	if docResult != nil {
		for _, s := range docResult.CodeSnippets {
			codeSnippets = append(codeSnippets, s.Code)
		}
		if len(docResult.InfoSnippets) > 0 {
			infoText = docResult.InfoSnippets[0].Content
		}
	}

	switch {
	case jsonFlag:
		output.PrintJSON(lib, infoText, codeSnippets)
	case importFlag:
		output.PrintImportOnly(lib)
	default:
		output.PrintResult(lib, infoText, codeSnippets)
		if !noCopyFlag {
			output.CopyImport(lib)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(searchCmd)
	rootCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output as JSON")
	rootCmd.Flags().BoolVar(&importFlag, "import-only", false, "Print only the import path")
	rootCmd.Flags().BoolVar(&noCopyFlag, "no-copy", false, "Skip clipboard copy")
	searchCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output as JSON")
	searchCmd.Flags().BoolVar(&importFlag, "import-only", false, "Print only the import path")
	searchCmd.Flags().BoolVar(&noCopyFlag, "no-copy", false, "Skip clipboard copy")
}

func resolve(cfg *config.Config, language, description string) ([]resolver.Result, error) {
	var llmCfg resolver.LLMConfig

	switch cfg.LLMProvider {
	case "gemini":
		if cfg.GeminiAPIKey != "" {
			llmCfg = resolver.LLMConfig{
				Provider: "gemini",
				APIKey:   cfg.GeminiAPIKey,
			}
		}
	}

	if llmCfg.APIKey != "" {
		r := resolver.NewLLM(llmCfg)
		results, err := r.Resolve(language, description)
		if err == nil && len(results) > 0 {
			return results, nil
		}
	}

	r := resolver.NewOffline()
	results, err := r.Resolve(language, description)
	if err != nil {
		return nil, fmt.Errorf("offline resolve: %w", err)
	}
	return results, nil
}
