package message

import (
	"encoding/json"
	"os"

	"github.com/summerwind/workflow-controller/pkg/webhook"
)

func Validate() error {
	msg := Message{}
	req := webhook.AdmissionRequest{
		Object: &msg,
	}

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.AdmissionResponse{
		UID: req.UID,
	}

	err = msg.Validate()
	res.Allowed = (err == nil)

	err = json.NewEncoder(os.Stdout).Encode(&res)
	if err != nil {
		return err
	}

	return nil
}
