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

var client *github.Client

func Reconcile() error {
	state := State{}
	err := json.NewDecoder(os.Stdin).Decode(&state)
	if err != nil {
		return err
	}

	if state.Resource.Status.URL == "" {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			return errors.New("GitHub access token must be specified")
		}

		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)

		r := state.Resource
		issue, _, err := client.Issues.Create(context.TODO(), r.Spec.Owner, r.Spec.Repository, &r.Spec.IssueRequest)
		if err != nil {
			return err
		}

		state.Resource.Status.URL = *issue.HTMLURL
		state.Resource.Status.Number = *issue.Number
		state.Resource.Status.CreationTime = time.Now().Unix()
	}

	err = json.NewEncoder(os.Stdout).Encode(&state)
	if err != nil {
		return err
	}

	return nil
}
