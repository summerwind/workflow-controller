package main

import (
	"github.com/spf13/cobra"

	"github.com/summerwind/workflow-controller/pkg/github/repository"
)

var repositoryValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return repository.Validate()
	},
}

func init() {
	repositoryCmd.AddCommand(repositoryValidateCmd)
}
