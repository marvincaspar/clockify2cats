package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	clockifyWorkspaceID          string
	clockifyUserID               string
	clockifyApiKey               string
	clockifyDescriptionDelimiter string
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize clockify2cats config",
		Long:  `Initialize clockify2cats config by providing workspace ID, user ID and api key.`,
		Run: func(cmd *cobra.Command, args []string) {
			viper.WriteConfig()
			viper.SafeWriteConfig()
		},
	}
}

func init() {
	initCmd := newInitCmd()
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")
	initCmd.PersistentFlags().StringVar(&clockifyWorkspaceID, "workspace", "", "Clockify workspace ID")
	initCmd.MarkPersistentFlagRequired("workspace")

	initCmd.PersistentFlags().StringVar(&clockifyUserID, "user", "", "Clockify user ID")
	initCmd.MarkPersistentFlagRequired("user")

	initCmd.PersistentFlags().StringVar(&clockifyApiKey, "api-key", "", "Clockify api key")
	initCmd.MarkPersistentFlagRequired("api-key")

	initCmd.PersistentFlags().StringVar(&clockifyDescriptionDelimiter, "description-delimiter", "#", "Clockify description delimiter to split description into text, text 2 and text external")

	viper.BindPFlag("workspace-id", initCmd.PersistentFlags().Lookup("workspace"))
	viper.BindPFlag("user-id", initCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("api-key", initCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("description-delimiter", initCmd.PersistentFlags().Lookup("description-delimiter"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
