package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/summerwind/workflow-controller/pkg/feed"
	"github.com/summerwind/workflow-controller/pkg/feed/v1alpha1"
)

type State struct {
	Object     *v1alpha1.Subscription `json:"object"`
	Dependents StateDependents        `json:"dependents"`
}

type StateDependents struct {
	Entries []*v1alpha1.Entry `json:"entry.v1alpha1.feed.workflow.summerwind.dev"`
}

var subscriptionReconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile resource",
	RunE:  subscriptionReconcile,
}

func init() {
	subscriptionCmd.AddCommand(subscriptionReconcileCmd)
}

func subscriptionReconcile(cmd *cobra.Command, args []string) error {
	state := State{}

	err := json.NewDecoder(os.Stdin).Decode(&state)
	if err != nil {
		return err
	}

	if !state.Object.NeedsUpdate() {
		return json.NewEncoder(os.Stdout).Encode(&state)
	}

	buf, err := fetch(state.Object.Spec.URL)
	if err != nil {
		return err
	}

	fetchTime := time.Now().Unix()
	if state.Object.Status.LastFetchTime == 0 {
		state.Object.Status.LastFetchTime = fetchTime
	}

	f, err := parse(buf)
	if err != nil {
		return err
	}

	// Workaround for behavior of XML decoder
	// See: https://github.com/golang/go/issues/8535
	if f.Channel != nil && f.Channel.Link == "" {
		f.Channel.Link = state.Object.Spec.URL
	}

	entries := newEntries(state.Object, f)
	for _, newEntry := range entries {
		var exists bool

		for _, entry := range state.Dependents.Entries {
			if newEntry.Spec.URL == entry.Spec.URL {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		state.Dependents.Entries = append(state.Dependents.Entries, newEntry)
	}

	return json.NewEncoder(os.Stdout).Encode(&state)
}

func fetch(feedURL string) ([]byte, error) {
	res, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func parse(data []byte) (*feed.Feed, error) {
	var f feed.Feed

	err := xml.Unmarshal(data, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func newEntries(sub *v1alpha1.Subscription, f *feed.Feed) []*v1alpha1.Entry {
	entries := []*v1alpha1.Entry{}

	// RSS Feed
	if f.Channel != nil {
		for _, item := range f.Channel.Items {
			t := time.Time(item.PubDate).Unix()

			if t < sub.Status.LastFetchTime {
				continue
			}

			e := v1alpha1.NewEntry()
			e.Name = fmt.Sprintf("%s-%d", sub.Name, t)
			e.Namespace = sub.Namespace
			e.Spec = v1alpha1.EntrySpec{
				URL:         item.Link,
				Title:       item.Title,
				PublishedAt: t,
				Feed: v1alpha1.EntrySpecFeed{
					Title: f.Channel.Title,
					URL:   f.Channel.Link,
				},
			}

			entries = append(entries, e)
		}

		return entries
	}

	// Atom Feed
	for _, fe := range f.Entries {
		t := time.Time(fe.Updated).Unix()

		if t < sub.Status.LastFetchTime {
			continue
		}

		e := v1alpha1.NewEntry()
		e.Name = fmt.Sprintf("%s-%d", sub.Name, t)
		e.Namespace = sub.Namespace
		e.Spec = v1alpha1.EntrySpec{
			URL:         fe.Link.Href,
			Title:       fe.Title,
			PublishedAt: t,
			Feed: v1alpha1.EntrySpecFeed{
				Title: f.Title,
				URL:   f.Link.Href,
			},
		}

		entries = append(entries, e)
	}

	return entries
}
