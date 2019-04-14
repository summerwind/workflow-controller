package repository

import (
	"errors"
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepositorySpec   `json:"spec,omitempty"`
	Status RepositoryStatus `json:"status,omitempty"`
}

func (r *Repository) Validate() error {
	return r.Spec.Validate()
}

type RepositorySpec struct {
	Owner         string            `json:"owner"`
	Description   *string           `json:"description,omitempty"`
	Homepage      *string           `json:"homepage,omitempty"`
	Private       *bool             `json:"private,omitempty"`
	HasIssues     *bool             `json:"hasIssues,omitempty"`
	HasProjects   *bool             `json:"hasProjects,omitempty"`
	HasWiki       *bool             `json:"hasWiki,omitempty"`
	DefaultBranch *string           `json:"defaultBranch,omitempty"`
	Topics        []string          `json:"topics,omitempty"`
	Labels        []RepositoryLabel `json:"labels,omitempty"`
}

func (r *RepositorySpec) Validate() error {
	if r.Owner == "" {
		return errors.New("owner must be specified")
	}

	labelNames := map[string]bool{}

	for i, label := range r.Labels {
		err := label.Validate()
		if err != nil {
			return fmt.Errorf("labels[%d]: %v", i, err)
		}

		_, ok := labelNames[label.Name]
		if ok {
			return fmt.Errorf("labels[%s]: duplicated name", i)
		}

		labelNames[label.Name] = true
	}

	return nil
}

type RepositoryLabel struct {
	Name        string  `json:"name"`
	Color       *string `json:"color,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (r *RepositoryLabel) Validate() error {
	if r.Name == "" {
		return errors.New("name must be specified")
	}

	if r.Color == nil {
		return errors.New("color must be specified")
	}

	return nil
}

type RepositoryStatus struct {
	URL            string `json:"url"`
	LastUpdateTime int64  `json:"lastUpdateTime"`
}

type State struct {
	Resource *Repository `json:"resource"`
}
