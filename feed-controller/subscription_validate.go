package main

import (
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/feed/subscription"
)

var subscriptionValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return subscription.Validate()
	},
}

func init() {
	subscriptionCmd.AddCommand(subscriptionValidateCmd)
}
