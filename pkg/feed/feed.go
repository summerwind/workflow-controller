package feed

import (
	"encoding/xml"
	"fmt"
	"time"
)

var timeLayouts []string

func init() {
	timeLayouts = []string{
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
		"2006-01-02T15:04:05+07:00",
	}
}

type Feed struct {
	// RSS
	Channel *Channel `xml:"channel"`

	// Atom
	Title   string  `xml:"title"`
	Link    Link    `xml:"link"`
	Entries []Entry `xml:"entry"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	PubDate FeedTime `xml:"pubDate"`
}

type Link struct {
	Type string `xml:"type,attr"`
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Entry struct {
	Title   string   `xml:"title"`
	Link    Link     `xml:"link"`
	Updated FeedTime `xml:"updated"`
}

type FeedTime time.Time

func (ft *FeedTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var (
		s      string
		t      time.Time
		parsed bool
	)

	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	for _, layout := range timeLayouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			parsed = true
			break
		}
	}

	if !parsed {
		return fmt.Errorf("Invalid time format: %s", s)
	}

	*ft = FeedTime(t)

	return nil
}
