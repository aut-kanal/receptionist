package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/kanalbot/receptionist"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Receptionist's version",
	Run: func(cmd *cobra.Command, args []string) {
		logVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func logVersion() {
	logrus.Info("version   > ", receptionist.Version)
	logrus.Info("buildtime > ", receptionist.BuildTime)
	logrus.Info("commit    > ", receptionist.Commit)
}
