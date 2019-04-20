package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/google/go-github/v24/github"
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/github/v1alpha1"
	"golang.org/x/oauth2"
)

type IssueState struct {
	Object *v1alpha1.Issue `json:"object"`
}

var issueReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE:  runIssueReconcile,
}

func init() {
	issueCmd.AddCommand(issueReconcileCmd)
}

func runIssueReconcile(cmd *cobra.Command, args []string) error {
	state := IssueState{}
	err := json.NewDecoder(os.Stdin).Decode(&state)
	if err != nil {
		return err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("GitHub access token must be specified")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	if state.Object.Status.URL == "" {
		o := state.Object
		issue, _, err := client.Issues.Create(context.TODO(), o.Spec.Owner, o.Spec.Repository, &o.Spec.IssueRequest)
		if err != nil {
			return err
		}

		now := time.Now().Unix()
		state.Object.Status.URL = *issue.HTMLURL
		state.Object.Status.Number = *issue.Number
		state.Object.Status.CreationTime = now
		state.Object.Status.LastUpdateTime = now
	} else if state.Object.NeedsUpdate() {
		o := state.Object
		_, _, err := client.Issues.Edit(context.TODO(), o.Spec.Owner, o.Spec.Repository, o.Status.Number, &o.Spec.IssueRequest)
		if err != nil {
			return err
		}

		state.Object.Status.LastUpdateTime = time.Now().Unix()
	}

	return json.NewEncoder(os.Stdout).Encode(&state)
}
