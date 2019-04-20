package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/slack/v1alpha1"
	"github.com/summerwind/workflow-controller/pkg/webhook"
)

var messageValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE:  runMessageValidate,
}

func init() {
	messageCmd.AddCommand(messageValidateCmd)
}

func runMessageValidate(cmd *cobra.Command, args []string) error {
	msg := v1alpha1.Message{}
	req := webhook.NewAdmissionRequest()

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.NewAdmissionResponse(req)
	err = req.GetObject(&msg)
	if err != nil {
		res.SetFailure(err.Error())
	} else {
		err = msg.Validate()
		if err != nil {
			res.SetFailure(err.Error())
		} else {
			res.SetSuccess()
		}
	}

	return json.NewEncoder(os.Stdout).Encode(&res)
}
