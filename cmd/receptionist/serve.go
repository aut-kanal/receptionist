package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/kanalbot/receptionist/telegram"
)

var (
	serveCmd = &cobra.Command{
		Use:   "start",
		Short: "Start bot",
		Run:   start,
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

func start(cmd *cobra.Command, args []string) {
	logVersion()

	telegram.InitBot()
}
