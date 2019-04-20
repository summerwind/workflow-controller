package v1alpha1

import (
	"errors"
	"time"

	"github.com/google/go-github/v24/github"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const issueTimeMergin = 10

type Issue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IssueSpec   `json:"spec,omitempty"`
	Status IssueStatus `json:"status,omitempty"`
}

func (i *Issue) Validate() error {
	return i.Spec.Validate()
}

func (i *Issue) NeedsUpdate() bool {
	if !i.Spec.Update {
		return false
	}

	if i.Status.LastUpdateTime+issueTimeMergin > time.Now().Unix() {
		return false
	}

	return true
}

type IssueSpec struct {
	github.IssueRequest
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	Update     bool   `json:"update"`
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
	URL            string `json:"url"`
	Number         int    `json:"number"`
	CreationTime   int64  `json:"creationTime"`
	LastUpdateTime int64  `json:"LastUpdateTime"`
}
