package main

import (
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/github/repository"
)

var repositoryReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return repository.Reconcile()
	},
}

func init() {
	repositoryCmd.AddCommand(repositoryReconcileCmd)
}
