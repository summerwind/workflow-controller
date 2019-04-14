package repository

import (
	"encoding/json"
	"os"

	"github.com/summerwind/workflow-controller/pkg/webhook"
)

func Validate() error {
	repo := Repository{}
	req := webhook.NewAdmissionRequest(repo)

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.NewAdmissionResponse(req)
	err = repo.Validate()
	if err != nil {
		res.SetFailure(err.Error())
	} else {
		res.SetSuccess()
	}

	err = json.NewEncoder(os.Stdout).Encode(&res)
	if err != nil {
		return err
	}

	return nil
}
