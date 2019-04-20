package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/slack/config"
	"github.com/summerwind/workflow-controller/pkg/slack/v1alpha1"
)

type State struct {
	Object *v1alpha1.Message `json:"object"`
}

var messageReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE:  runMessageReconcile,
}

func init() {
	messageCmd.AddCommand(messageReconcileCmd)
}

func runMessageReconcile(cmd *cobra.Command, args []string) error {
	state := State{}
	err := json.NewDecoder(os.Stdin).Decode(&state)
	if err != nil {
		return err
	}

	if state.Object.Status.SendTime != 0 {
		return json.NewEncoder(os.Stdout).Encode(&state)
	}

	log.Print("Loading config...")
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

	log.Print("Sending message...")
	res, err := http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %v", res.Status)
	}

	state.Object.Status.SendTime = time.Now().Unix()

	return json.NewEncoder(os.Stdout).Encode(&state)
}
