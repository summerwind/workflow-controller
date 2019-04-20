package v1alpha1

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Run struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RunSpec   `json:"spec"`
	Status RunStatus `json:"status"`
}

func (r *Run) Validate() error {
	return r.Spec.Validate()
}

type RunSpec struct {
	Source    *RunSpecSource    `json:"source"`
	Workspace string            `json:"workspace"`
	Vars      map[string]string `json:"vars,omitempty"`
}

func (r *RunSpec) Validate() error {
	if r.Source == nil {
		return errors.New("source must be specified")
	}

	err := r.Source.Validate()
	if err != nil {
		return err
	}

	return nil
}

type RunSpecSource struct {
	Git  *RunSpecSourceGit  `json:"git"`
	File *RunSpecSourceFile `json:"file"`
}

func (r *RunSpecSource) Validate() error {
	if r.Git != nil {
		err := r.Git.Validate()
		if err != nil {
			return err
		}
	}

	if r.File != nil {
		err := r.File.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type RunSpecSourceGit struct {
	URL  string `json:"url"`
	Ref  string `json:"ref"`
	Path string `json:"path"`
}

func (r *RunSpecSourceGit) Validate() error {
	if r.URL == "" {
		return errors.New("url must be specified")
	}

	if r.Path == "" {
		return errors.New("path must be specified")
	}

	return nil
}

type RunSpecSourceFile struct {
	Path string `json:"path"`
}

func (r *RunSpecSourceFile) Validate() error {
	if r.Path == "" {
		return errors.New("path must be specified")
	}

	return nil
}

type RunStatus struct {
	LastApplyTime   int64  `json:"lastApplyTime"`
	LastApplyCommit string `json:"lastApplyCommit"`
}

type State struct {
	Resource *Run `json:"resource"`
}
