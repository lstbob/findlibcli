package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "findlib",
	Short: "Find libraries by natural language description",
	Long: `findlib searches for libraries and APIs matching a description in a given language.

Examples:
  findlib go "http client with retries"
  findlib python "json parsing"
  findlib js "state management"
  findlib config gemini "your-api-key"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
