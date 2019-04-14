package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func Reconcile() error {
	state := State{}
	err := json.NewDecoder(os.Stdin).Decode(&state)
	if err != nil {
		return err
	}

	if state.Resource.Status.SendTime == 0 {
		channels := getChannels()
		url, ok := channels[state.Resource.Spec.Channel]
		if !ok {
			return fmt.Errorf("invalid channel name: %v", state.Resource.Spec.Channel)
		}
		if url == "" {
			return fmt.Errorf("invalid channel URL: %v", state.Resource.Spec.Channel)
		}

		buf, err := json.Marshal(state.Resource.Spec)
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

		state.Resource.Status.SendTime = time.Now().Unix()
	}

	err = json.NewEncoder(os.Stdout).Encode(&state)
	if err != nil {
		return err
	}

	return nil
}

func getChannels() map[string]string {
	channels := map[string]string{}

	chStr := os.Getenv("SLACK_CHANNELS")
	chunks := strings.Split(chStr, ",")
	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		nameAndURL := strings.Split(chunk, "=")
		channels[nameAndURL[0]] = nameAndURL[1]
	}

	return channels
}
