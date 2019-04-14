package run

import (
	"encoding/json"
	"os"

	"github.com/summerwind/workflow-controller/pkg/webhook"
)

func Validate() error {
	run := Run{}
	req := webhook.AdmissionRequest{
		Object: &run,
	}

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.AdmissionResponse{
		UID: req.UID,
	}

	err = run.Validate()
	res.Allowed = (err == nil)

	err = json.NewEncoder(os.Stdout).Encode(&res)
	if err != nil {
		return err
	}

	return nil
}
