package main

import (
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/github/issue"
)

var issueReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return issue.Reconcile()
	},
}

func init() {
	issueCmd.AddCommand(issueReconcileCmd)
}
