package main

import (
	"github.com/spf13/cobra"
)

var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Manage Repository resource",
}

func init() {
	rootCmd.AddCommand(repositoryCmd)
}
