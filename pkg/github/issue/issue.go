package issue

import (
	"errors"

	"github.com/google/go-github/v24/github"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Issue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IssueSpec   `json:"spec,omitempty"`
	Status IssueStatus `json:"status,omitempty"`
}

func (i *Issue) Validate() error {
	return i.Spec.Validate()
}

type IssueSpec struct {
	github.IssueRequest
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
}

func (i *IssueSpec) Validate() error {
	if i.Owner == "" {
		return errors.New("owner must be specified")
	}

	if *i.IssueRequest.Title == "" {
		return errors.New("title must be specified")
	}

	return nil
}

type IssueStatus struct {
	URL          string `json:"url"`
	Number       int    `json:"number"`
	CreationTime int64  `json:"creationTime"`
}

type State struct {
	Resource *Issue `json:"resource"`
}
