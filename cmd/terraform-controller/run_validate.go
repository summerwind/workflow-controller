package main

import (
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/terraform/run"
)

var runValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.Validate()
	},
}

func init() {
	runCmd.AddCommand(runValidateCmd)
}
