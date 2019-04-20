package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/summerwind/workflow-controller/pkg/slack/config"
)

type State struct {
	Object *Message `json:"object"`
}

func Reconcile() error {
	state := State{}
	err := json.NewDecoder(os.Stdin).Decode(&state)
	if err != nil {
		return err
	}

	if state.Object.Status.SendTime == 0 {
		c, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		url, ok := c.Channels[state.Object.Spec.Channel]
		if !ok {
			return fmt.Errorf("invalid channel name: %v", state.Object.Spec.Channel)
		}
		if url == "" {
			return fmt.Errorf("invalid channel URL: %v", state.Object.Spec.Channel)
		}

		buf, err := json.Marshal(state.Object.WebhookMessage())
		if err != nil {
			return err
		}

		res, err := http.Post(url, "application/json", bytes.NewReader(buf))
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected response status: %v", res.Status)
		}

		state.Object.Status.SendTime = time.Now().Unix()
	}

	err = json.NewEncoder(os.Stdout).Encode(&state)
	if err != nil {
		return err
	}

	return nil
}
