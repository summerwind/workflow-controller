package main

import (
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/github/issue"
)

var issueValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return issue.Validate()
	},
}

func init() {
	issueCmd.AddCommand(issueValidateCmd)
}