package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/github/v1alpha1"
	"github.com/summerwind/workflow-controller/pkg/webhook"
)

var repositoryValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE:  runRepositoryValidate,
}

func init() {
	repositoryCmd.AddCommand(repositoryValidateCmd)
}

func runRepositoryValidate(cmd *cobra.Command, args []string) error {
	repo := v1alpha1.Repository{}
	req := webhook.NewAdmissionRequest()

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.NewAdmissionResponse(req)
	err = req.GetObject(&repo)
	if err != nil {
		res.SetFailure(err.Error())
	} else {
		err = repo.Validate()
		if err != nil {
			res.SetFailure(err.Error())
		} else {
			res.SetSuccess()
		}
	}

	return json.NewEncoder(os.Stdout).Encode(&res)
}
