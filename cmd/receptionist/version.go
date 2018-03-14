package main

import (
	"git.cafebazaar.ir/arya/baargir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	logrus.Info("version   > ", baargir.Version)
	logrus.Info("buildtime > ", baargir.BuildTime)
	logrus.Info("commit    > ", baargir.Commit)
}
