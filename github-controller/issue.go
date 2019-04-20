package main

import (
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage Issue resource",
}

func init() {
	rootCmd.AddCommand(issueCmd)
}
