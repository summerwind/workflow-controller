package v1alpha1

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Entry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntrySpec   `json:"spec,omitempty"`
	Status EntryStatus `json:"status,omitempty"`
}

func NewEntry() *Entry {
	return &Entry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "feed.workflow.summerwind.dev/v1alpha1",
			Kind:       "Entry",
		},
		ObjectMeta: metav1.ObjectMeta{},
	}
}

func (e *Entry) Validate() error {
	return e.Spec.Validate()
}

type EntrySpec struct {
	URL         string        `json:"url"`
	Title       string        `json:"title"`
	PublishedAt int64         `json:"publishedAt"`
	Feed        EntrySpecFeed `json:"feed"`
}

func (e *EntrySpec) Validate() error {
	if e.Title == "" {
		return errors.New("title must be specified")
	}

	if e.URL == "" {
		return errors.New("url must be specified")
	}

	err := e.Feed.Validate()
	if err != nil {
		return fmt.Errorf("feed: %v", err)
	}

	return nil
}

type EntrySpecFeed struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func (e *EntrySpecFeed) Validate() error {
	if e.Title == "" {
		return errors.New("title must be specified")
	}

	if e.URL == "" {
		return errors.New("url must be specified")
	}

	return nil
}

type EntryStatus struct {
}
