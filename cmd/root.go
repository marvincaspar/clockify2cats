package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "clockify2cats",
	Short: "Tool to convert clockify time entries to cats time entries",
	Long:  `This tool allows you to convert clockify time entries to cats time entries.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Find config directory.
	configDir, err := os.UserConfigDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".clockify2cats" (without extension).
	viper.AddConfigPath(configDir + "/clockify2cats")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	viper.ReadInConfig()
}
