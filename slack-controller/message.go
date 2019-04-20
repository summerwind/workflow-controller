package main

import (
	"github.com/spf13/cobra"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Manage Message resource",
}

func init() {
	rootCmd.AddCommand(messageCmd)
}
