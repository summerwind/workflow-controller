package subscription

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/summerwind/workflow-controller/pkg/feed/entry"
)

var fetchMergin int64 = 30

type State struct {
	Object     *Subscription   `json:"object"`
	Dependents StateDependents `json:"dependents"`
}

type StateDependents struct {
	Entries []*entry.Entry `json:"entry.v1alpha1.feed.workflow.summerwind.dev"`
}

func Reconcile() error {
	state := State{}

	err := json.NewDecoder(os.Stdin).Decode(&state)
	if err != nil {
		return err
	}

	buf, err := fetch(state.Object.Spec.URL)
	if err != nil {
		return err
	}

	fetchedTime := time.Now().Unix()
	if state.Object.Status.LastFetchedTime == 0 {
		state.Object.Status.LastFetchedTime = fetchedTime
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

	// Prevent time update loolping by using
	if state.Object.Status.LastFetchedTime+fetchMergin < fetchedTime {
		state.Object.Status.LastFetchedTime = fetchedTime
	}

	err = json.NewEncoder(os.Stdout).Encode(&state)
	if err != nil {
		return err
	}

	return nil
}

func fetch(feedURL string) ([]byte, error) {
	res, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func parse(data []byte) (*Feed, error) {
	var feed Feed

	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func newEntries(sub *Subscription, f *Feed) []*entry.Entry {
	entries := []*entry.Entry{}

	// RSS Feed
	if f.Channel != nil {
		for _, item := range f.Channel.Items {
			t := time.Time(item.PubDate).Unix()

			if t < sub.Status.LastFetchedTime {
				continue
			}

			e := entry.NewEntry()
			e.Name = fmt.Sprintf("%s-%d", sub.Name, t)
			e.Namespace = sub.Namespace
			e.Spec = entry.EntrySpec{
				URL:         item.Link,
				Title:       item.Title,
				PublishedAt: t,
				Feed: entry.EntrySpecFeed{
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

		if t < sub.Status.LastFetchedTime {
			continue
		}

		e := entry.NewEntry()
		e.Name = fmt.Sprintf("%s-%d", sub.Name, t)
		e.Namespace = sub.Namespace
		e.Spec = entry.EntrySpec{
			URL:         fe.Link.Href,
			Title:       fe.Title,
			PublishedAt: t,
			Feed: entry.EntrySpecFeed{
				Title: f.Title,
				URL:   f.Link.Href,
			},
		}

		entries = append(entries, e)
	}

	return entries
}
