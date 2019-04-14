package run

import (
	"errors"
	"log"
	"os"

	"github.com/summerwind/workflow-controller/pkg/webhook"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

type Run struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Metadata   interface{} `json:"metadata"`
	Spec       RunSpec     `json:"spec"`
	Status     RunStatus   `json:"status"`
}

func (r *Run) Validate() error {
	err := r.Spec.Validate()
	if err != nil {
		return err
	}

	return nil
}

type RunSpec struct {
	Source    RunSpecSource     `json:"source"`
	Workspace string            `json:"workspace"`
	Vars      map[string]string `json:"vars,omitempty"`
}

func (r *RunSpec) Validate() error {
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

type AdmissionRequest struct {
	webhook.AdmissionRequest
	Object *Run `json:"object"`
}
