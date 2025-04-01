package cmd

import "github.com/spf13/cobra"

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Find a single struct in a file",
	Long:  "Locate and display a specific struct definition within a given source file",
	RunE:  structCmdRun,
}

func init() {
	rootCmd.AddCommand(structCmd)
}

func structCmdRun(cmd *cobra.Command, args []string) error {
	return nil
}
