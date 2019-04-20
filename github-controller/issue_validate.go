package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/github/v1alpha1"
	"github.com/summerwind/workflow-controller/pkg/webhook"
)

var issueValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE:  runIssueValidate,
}

func init() {
	issueCmd.AddCommand(issueValidateCmd)
}

func runIssueValidate(cmd *cobra.Command, args []string) error {
	issue := v1alpha1.Issue{}
	req := webhook.NewAdmissionRequest()

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.NewAdmissionResponse(req)
	err = req.GetObject(&issue)
	if err != nil {
		res.SetFailure(err.Error())
	} else {
		err = issue.Validate()
		if err != nil {
			res.SetFailure(err.Error())
		} else {
			res.SetSuccess()
		}
	}

	return json.NewEncoder(os.Stdout).Encode(&res)
}
