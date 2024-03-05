package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "imgo",
	Short: "improve clarity of photos through an AI model",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
