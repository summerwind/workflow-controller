package main

import (
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/terraform/run"
)

var runReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.Reconcile()
	},
}

func init() {
	runCmd.AddCommand(runReconcileCmd)
}
