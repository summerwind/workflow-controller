package main

import (
	"github.com/spf13/cobra"
)

var subscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "Manage Subscription resource",
}

func init() {
	rootCmd.AddCommand(subscriptionCmd)
}
