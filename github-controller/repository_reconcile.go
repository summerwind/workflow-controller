package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v24/github"
	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/github/v1alpha1"
	"golang.org/x/oauth2"
)

type RepositoryState struct {
	Object *v1alpha1.Repository `json:"object"`
}

var repositoryReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE:  runRepositoryReconcile,
}

func init() {
	repositoryCmd.AddCommand(repositoryReconcileCmd)
}

func runRepositoryReconcile(cmd *cobra.Command, args []string) error {
	state := RepositoryState{}
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
		var (
			repo *github.Repository
			res  *github.Response
			err  error
		)

		repo, res, err = getRepository(client, state.Object)
		if err != nil {
			return err
		}

		if res.StatusCode == 404 {
			repo, res, err = createRepository(client, state.Object)
			if err != nil {
				return err
			}
		}

		now := time.Now().Unix()
		state.Object.Status.URL = repo.GetHTMLURL()
		state.Object.Status.CreationTime = now
	}

	if state.Object.NeedsUpdate() {
		repo, res, err := updateRepository(client, state.Object)
		if err != nil {
			return err
		}

		if res.StatusCode/100 == 2 {
			err = updateTopics(client, state.Object)
			if err != nil {
				return err
			}

			err = updateLabels(client, state.Object)
			if err != nil {
				return err
			}

			state.Object.Status.URL = repo.GetHTMLURL()
			state.Object.Status.LastUpdateTime = time.Now().Unix()
		} else if res.StatusCode == 404 {
			// The repository seems to be deleted and needs to be recreated.
			state.Object.Status.URL = ""
		}
	}

	return json.NewEncoder(os.Stdout).Encode(&state)
}

func getRepository(client *github.Client, repo *v1alpha1.Repository) (*github.Repository, *github.Response, error) {
	return client.Repositories.Get(context.TODO(), repo.Spec.Owner, repo.Name)
}

func createRepository(client *github.Client, repo *v1alpha1.Repository) (*github.Repository, *github.Response, error) {
	r := github.Repository{
		Name:        &repo.Name,
		Description: &repo.Spec.Description,
		Homepage:    &repo.Spec.Homepage,
		Private:     &repo.Spec.Private,
		HasIssues:   &repo.Spec.HasIssues,
		HasProjects: &repo.Spec.HasProjects,
		HasWiki:     &repo.Spec.HasWiki,
	}

	log.Print("Creating repository")
	return client.Repositories.Create(context.TODO(), repo.Spec.Owner, &r)
}

func updateRepository(client *github.Client, repo *v1alpha1.Repository) (*github.Repository, *github.Response, error) {
	r := github.Repository{
		Name:        &repo.Name,
		Description: &repo.Spec.Description,
		Homepage:    &repo.Spec.Homepage,
		Private:     &repo.Spec.Private,
		HasIssues:   &repo.Spec.HasIssues,
		HasProjects: &repo.Spec.HasProjects,
		HasWiki:     &repo.Spec.HasWiki,
	}

	if repo.Spec.DefaultBranch != "" {
		r.DefaultBranch = &repo.Spec.DefaultBranch
	}

	log.Print("Updating repository")
	return client.Repositories.Edit(context.TODO(), repo.Spec.Owner, repo.Name, &r)
}

func updateTopics(client *github.Client, repo *v1alpha1.Repository) error {
	log.Print("Updating topics")
	_, _, err := client.Repositories.ReplaceAllTopics(context.TODO(), repo.Spec.Owner, repo.Name, repo.Spec.Topics)
	return err
}

func updateLabels(client *github.Client, repo *v1alpha1.Repository) error {
	ledger := map[string]v1alpha1.RepositoryLabel{}
	for _, label := range repo.Spec.Labels {
		ledger[label.Name] = label
	}

	updateLabels := []v1alpha1.RepositoryLabel{}
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

			if label.Color != l.GetColor() || label.Description != l.GetDescription() {
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
			Color:       &l.Color,
			Description: &l.Description,
		}

		log.Printf("Creating label %s", l.Name)
		_, _, err := client.Issues.CreateLabel(context.TODO(), repo.Spec.Owner, repo.Name, &label)
		if err != nil {
			return err
		}
	}

	for _, l := range updateLabels {
		label := github.Label{
			Name:        &l.Name,
			Color:       &l.Color,
			Description: &l.Description,
		}

		log.Printf("Updating label %s", l.Name)
		_, _, err := client.Issues.EditLabel(context.TODO(), repo.Spec.Owner, repo.Name, l.Name, &label)
		if err != nil {
			return err
		}
	}

	for _, name := range deleteLabels {
		log.Printf("Deleting label %s", name)
		_, err := client.Issues.DeleteLabel(context.TODO(), repo.Spec.Owner, repo.Name, name)
		if err != nil {
			return err
		}
	}

	return nil
}
