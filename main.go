package main

import (
	"fmt"
	"os"

	"css-color-variable-creator/cmd/create"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "css-color-variable-creator",
	Short: "A CLI tool to manage CSS color variables",
	Long: `css-color-variable-creator is a command line tool that helps you
create and manage CSS color variables in your stylesheets.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This is the default command when no subcommand is provided
		fmt.Println("Welcome to CSS Color Variable Creator!")
		fmt.Println("Use --help to see available commands")
	},
}

func init() {
	rootCmd.AddCommand(create.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
