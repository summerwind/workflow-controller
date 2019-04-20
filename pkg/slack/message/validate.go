package message

import (
	"encoding/json"
	"os"

	"github.com/summerwind/workflow-controller/pkg/slack/config"
	"github.com/summerwind/workflow-controller/pkg/webhook"
)

func Validate() error {
	msg := Message{}
	req := webhook.NewAdmissionRequest(&msg)

	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		return err
	}

	res := webhook.NewAdmissionResponse(req)
	err = msg.Validate()
	if err != nil {
		res.SetFailure(err.Error())
	} else {
		c, err := config.Load()
		if err != nil {
			res.SetFailure("internal error: failed to load config")
		}

		url, ok := c.Channels[msg.Spec.Channel]
		if !ok {
			res.SetFailure("invalid channel name")
		} else if url == "" {
			res.SetFailure("invalid channel URL")
		} else {
			res.SetSuccess()
		}
	}

	err = json.NewEncoder(os.Stdout).Encode(&res)
	if err != nil {
		return err
	}

	return nil
}
