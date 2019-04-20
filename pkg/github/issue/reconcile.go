package issue

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
)

type State struct {
	Object *Issue `json:"object"`
}

func Reconcile() error {
	state := State{}
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

	err = json.NewEncoder(os.Stdout).Encode(&state)
	if err != nil {
		return err
	}

	return nil
}
