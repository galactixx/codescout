package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codescout",
	Short: "Quickly locate classes and functions in any codebase",
	Long: `CodeScout is a lightweight CLI tool and Go module that helps developers 
quickly locate classes, functions, and other elements in any codebase. 
It makes code navigation and structure analysis fast and efficient.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
