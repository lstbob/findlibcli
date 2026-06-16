package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
	"github.com/lstbob/findlibcli/internal/resolver"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF87")).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00BFFF"))

	importStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Italic(true)

	codeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			PaddingLeft(2)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0"))
)

type JSONOutput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ImportPath  string   `json:"import_path"`
	Language    string   `json:"language"`
	CodeSnippets []string `json:"code_snippets,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

func PrintResult(lib *resolver.Result, docs string, codeSnippets []string) {
	isTTY := isTerminal()

	if isTTY {
		printStyled(lib, docs, codeSnippets)
	} else {
		printPlain(lib, docs, codeSnippets)
	}
}

func PrintJSON(lib *resolver.Result, docs string, codeSnippets []string) {
	out := JSONOutput{
		Name:        lib.Name,
		Description: lib.Description,
		ImportPath:  lib.ImportPath,
		Language:    lib.Language,
		CodeSnippets: codeSnippets,
		Documentation: docs,
	}
	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(data))
}

func PrintImportOnly(lib *resolver.Result) {
	fmt.Println(lib.ImportPath)
}

func CopyImport(lib *resolver.Result) {
	if err := clipboard.WriteAll(lib.ImportPath); err == nil {
		fmt.Println("✓ Import path copied to clipboard")
	}
}

func isTerminal() bool {
	return true
}

func printStyled(lib *resolver.Result, docs string, codeSnippets []string) {
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(titleStyle.Render(lib.Name))
	fmt.Println(sectionStyle.Render("Description:"))
	fmt.Println(codeStyle.Render(lib.Description))
	fmt.Println()
	fmt.Println(sectionStyle.Render("Import / Install:"))
	fmt.Println(importStyle.Render("  " + lib.ImportPath))
	fmt.Println()

	if len(codeSnippets) > 0 {
		fmt.Println(sectionStyle.Render("Usage:"))
		for _, s := range codeSnippets {
			for _, line := range strings.Split(s, "\n") {
				fmt.Println(codeStyle.Render(line))
			}
		}
		fmt.Println()
	}

	if docs != "" {
		fmt.Println(sectionStyle.Render("Documentation:"))
		fmt.Println(infoStyle.Render("  " + docs))
		fmt.Println()
	}

	fmt.Println(strings.Repeat("─", 50))
}

func printPlain(lib *resolver.Result, docs string, codeSnippets []string) {
	fmt.Printf("Library: %s\n", lib.Name)
	fmt.Printf("Description: %s\n", lib.Description)
	fmt.Printf("Import: %s\n", lib.ImportPath)
	if len(codeSnippets) > 0 {
		fmt.Printf("Usage:\n")
		for _, s := range codeSnippets {
			fmt.Println(s)
		}
	}
	if docs != "" {
		fmt.Printf("Docs: %s\n", docs)
	}
}
