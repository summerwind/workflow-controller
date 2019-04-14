package repository

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

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("GitHub access token must be specified")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)

	url, err := updateRepository(state.Resource)
	if err != nil {
		return err
	}

	err = updateTopics(state.Resource)
	if err != nil {
		return err
	}

	err = updateLabels(state.Resource)
	if err != nil {
		return err
	}

	state.Resource.Status.URL = url
	state.Resource.Status.LastUpdateTime = time.Now().Unix()

	err = json.NewEncoder(os.Stdout).Encode(&state)
	if err != nil {
		return err
	}

	return nil
}

func updateRepository(repo *Repository) (string, error) {
	r := github.Repository{
		Name:          &repo.Name,
		Description:   repo.Spec.Description,
		Homepage:      repo.Spec.Homepage,
		Private:       repo.Spec.Private,
		HasIssues:     repo.Spec.HasIssues,
		HasProjects:   repo.Spec.HasProjects,
		HasWiki:       repo.Spec.HasWiki,
		DefaultBranch: repo.Spec.DefaultBranch,
	}

	logger.Print("Updating repository")
	updated, _, err := client.Repositories.Edit(context.TODO(), repo.Spec.Owner, repo.Name, &r)

	return *updated.HTMLURL, err
}

func updateTopics(repo *Repository) error {
	logger.Print("Updating topics")
	_, _, err := client.Repositories.ReplaceAllTopics(context.TODO(), repo.Spec.Owner, repo.Name, repo.Spec.Topics)
	return err
}

func updateLabels(repo *Repository) error {
	ledger := map[string]RepositoryLabel{}
	for _, label := range repo.Spec.Labels {
		ledger[label.Name] = label
	}

	updateLabels := []RepositoryLabel{}
	deleteLabels := []string{}

	page := 1
	for {
		opts := github.ListOptions{
			Page:    page,
			PerPage: 100,
		}

		labels, res, err := client.Issues.ListLabels(context.TODO(), repo.Spec.Owner, repo.Name, &opts)
		if err != nil {
			return err
		}

		for _, l := range labels {
			label, ok := ledger[*l.Name]
			if !ok {
				deleteLabels = append(deleteLabels, *l.Name)
				continue
			}

			if *label.Color != *l.Color || *label.Description != *l.Description {
				updateLabels = append(updateLabels, label)
			}

			delete(ledger, *l.Name)
		}

		if res.LastPage == 0 {
			break
		}

		page = res.NextPage
	}

	for _, l := range ledger {
		label := github.Label{
			Name:        &l.Name,
			Color:       l.Color,
			Description: l.Description,
		}

		logger.Printf("Creating label %s", l.Name)
		_, _, err := client.Issues.CreateLabel(context.TODO(), repo.Spec.Owner, repo.Name, &label)
		if err != nil {
			return err
		}
	}

	for _, l := range updateLabels {
		label := github.Label{
			Name:        &l.Name,
			Color:       l.Color,
			Description: l.Description,
		}

		logger.Printf("Updating label %s", l.Name)
		_, _, err := client.Issues.EditLabel(context.TODO(), repo.Spec.Owner, repo.Name, l.Name, &label)
		if err != nil {
			return err
		}
	}

	for _, name := range deleteLabels {
		logger.Printf("Deleting label %s", name)
		_, err := client.Issues.DeleteLabel(context.TODO(), repo.Spec.Owner, repo.Name, name)
		if err != nil {
			return err
		}
	}

	return nil
}
