package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/feed/v1alpha1"
	"github.com/summerwind/workflow-controller/pkg/webhook"
)

var subscriptionValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resource",
	RunE:  runSubscriptionValidate,
}

func init() {
	subscriptionCmd.AddCommand(subscriptionValidateCmd)
}

func runSubscriptionValidate(cmd *cobra.Command, args []string) error {
	sub := v1alpha1.Subscription{}
	req := webhook.NewAdmissionRequest()

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.NewAdmissionResponse(req)
	err = req.GetObject(&sub)
	if err != nil {
		res.SetFailure(err.Error())
	} else {
		err = sub.Validate()
		if err != nil {
			res.SetFailure(err.Error())
		} else {
			res.SetSuccess()
		}
	}

	return json.NewEncoder(os.Stdout).Encode(&res)
}
