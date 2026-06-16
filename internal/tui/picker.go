package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/lstbob/findlibcli/internal/resolver"
)

func PickLibrary(results []resolver.Result) (*resolver.Result, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("no results to pick from")
	}
	if len(results) == 1 {
		return &results[0], nil
	}

	if !isTerminal() {
		return pickFromStdin(results)
	}

	opts := make([]huh.Option[int], len(results))
	for i, r := range results {
		label := fmt.Sprintf("%s  (%s)", r.Name, r.ImportPath)
		if r.Description != "" {
			label = fmt.Sprintf("%s  —  %s", r.Name, truncate(r.Description, 60))
		}
		opts[i] = huh.Option[int]{Key: label, Value: i}
	}

	var selected int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select a library").
				Options(opts...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("picker: %w", err)
	}

	return &results[selected], nil
}

func isTerminal() bool {
	info, err := os.Stdout.Stat()
	return err == nil && (info.Mode()&os.ModeCharDevice) != 0
}

func pickFromStdin(results []resolver.Result) (*resolver.Result, error) {
	fmt.Println("\nMultiple libraries found. Pipe to a TTY for interactive selection.")
	for i, r := range results {
		fmt.Printf("[%d] %s\n    %s\n    import: %s\n\n", i+1, r.Name, r.Description, r.ImportPath)
	}
	fmt.Print("Enter number (1-" + fmt.Sprint(len(results)) + "): ")
	var n int
	_, err := fmt.Scanf("%d", &n)
	if err != nil || n < 1 || n > len(results) {
		return nil, fmt.Errorf("invalid selection")
	}
	return &results[n-1], nil
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}
