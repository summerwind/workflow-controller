package main

import (
	"github.com/spf13/cobra"

	"github.com/summerwind/workflow-controller/pkg/slack/message"
)

var messageReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return message.Reconcile()
	},
}

func init() {
	messageCmd.AddCommand(messageReconcileCmd)
}
