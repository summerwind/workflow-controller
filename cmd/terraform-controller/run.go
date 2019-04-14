package main

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage Run resource",
}

func init() {
	rootCmd.AddCommand(runCmd)
}
