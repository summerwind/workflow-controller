package main

import (
	"github.com/spf13/cobra"

	"github.com/summerwind/workflow-controller/pkg/slack/message"
)

var messageValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return message.Validate()
	},
}

func init() {
	messageCmd.AddCommand(messageValidateCmd)
}
