package cmd

import (
	"github.com/spf13/cobra"
)

var Version = "0.0.1" // Default version number

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Print the version number of the application",
		Run: func(cmd *cobra.Command, args []string) {
			// Print the version number
			cmd.Printf("Version: %s\n", Version)
		},
	}
}

func init() {
	// Add the version command to the root command
	rootCmd.AddCommand(newVersionCmd())
}
