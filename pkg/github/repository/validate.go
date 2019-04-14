package repository

import (
	"encoding/json"
	"os"

	"github.com/summerwind/workflow-controller/pkg/webhook"
)

func Validate() error {
	repo := Repository{}
	req := webhook.AdmissionRequest{
		Object: &repo,
	}

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.AdmissionResponse{
		UID: req.UID,
	}

	err = repo.Validate()
	res.Allowed = (err == nil)

	err = json.NewEncoder(os.Stdout).Encode(&res)
	if err != nil {
		return err
	}

	return nil
}
