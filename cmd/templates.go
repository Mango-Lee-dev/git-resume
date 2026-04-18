package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wootaiklee/git-resume/internal/llm"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available prompt templates",
	Long: `Display all available prompt templates that can be used with the analyze command.

Templates customize the tone and focus of generated resume bullet points
for different industries, roles, or company cultures.

Example:
  git-resume templates
  git-resume analyze --template=startup`,
	Run: runTemplates,
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}

func runTemplates(cmd *cobra.Command, args []string) {
	fmt.Println("Available Templates")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	for name, tmpl := range llm.BuiltinTemplates {
		fmt.Printf("  %s\n", name)
		fmt.Printf("    Name: %s\n", tmpl.Name)
		fmt.Printf("    Description: %s\n", tmpl.Description)
		fmt.Printf("    Tone: %s\n", tmpl.ToneStyle)
		fmt.Printf("    Focus: ")
		for i, f := range tmpl.Focus {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", f)
		}
		fmt.Println()
		fmt.Println()
	}

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("Usage: git-resume analyze --template=<name>")
	fmt.Println()
	fmt.Println("Custom templates can be loaded from JSON files.")
	fmt.Println("See documentation for template file format.")
}
