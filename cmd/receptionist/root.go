package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/kanalbot/receptionist/configuration"
)

var rootCmd = &cobra.Command{
	Use:   "receptionistd <subcommand>",
	Short: "Kanal's bot",
	Long:  "Receptionist receives user's message",
	Run:   nil,
}

func init() {
	cobra.OnInitialize()

	configFilePath := "config.yaml"
	rootCmd.PersistentFlags().StringVarP(&configFilePath,
		"config", "c", configFilePath, "Path to the config file (eg ./config.yaml)")
	configuration.SetFilePath(configFilePath)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if configuration.GetInstance().GetBool("debug") {
			configuration.SetDebugLogLevel(true)
		}
	}
}
