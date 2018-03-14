package main

import (
	"git.cafebazaar.ir/arya/baargir/configuration"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "receptionistd <subcommand>",
	Short: "Kanal's bot",
	Long:  "Receptionist receives user's message",
	Run:   nil,
}

func init() {
	var isDebug bool
	configFilePath := "config.yaml"

	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringVarP(&configFilePath,
		"Config-file", "c", configFilePath, "Path to the config file (eg ./config.toml)")

	configuration.SetFilePath(configFilePath)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if isDebug {
			configuration.GetInstance().Set("debug", true)
			configuration.SetDebugLogLevel(true)
		}
	}
}
