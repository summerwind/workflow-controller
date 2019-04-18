package main

import (
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/feed/subscription"
)

var subscriptionReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return subscription.Reconcile()
	},
}

func init() {
	subscriptionCmd.AddCommand(subscriptionReconcileCmd)
}
