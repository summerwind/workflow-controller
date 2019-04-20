package message

import (
	"encoding/json"
	"errors"

	"github.com/nlopes/slack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Message struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MessageSpec   `json:"spec,omitempty"`
	Status MessageStatus `json:"status,omitempty"`
}

func (m *Message) Validate() error {
	return m.Spec.Validate()
}

func (m *Message) WebhookMessage() *slack.WebhookMessage {
	attachments := []slack.Attachment{}
	for _, a := range m.Spec.Attachments {
		fields := []slack.AttachmentField{}
		for _, f := range a.Fields {
			fields = append(fields, slack.AttachmentField(f))
		}

		actions := []slack.AttachmentAction{}
		for _, ac := range a.Actions {
			options := []slack.AttachmentActionOption{}
			for _, o := range ac.Options {
				options = append(options, slack.AttachmentActionOption(o))
			}

			selectedOptions := []slack.AttachmentActionOption{}
			for _, o := range ac.SelectedOptions {
				selectedOptions = append(selectedOptions, slack.AttachmentActionOption(o))
			}

			optionGroups := []slack.AttachmentActionOptionGroup{}
			for _, og := range ac.OptionGroups {
				goptions := []slack.AttachmentActionOption{}
				for _, o := range og.Options {
					goptions = append(goptions, slack.AttachmentActionOption(o))
				}

				optionGroups = append(optionGroups, slack.AttachmentActionOptionGroup{
					Text:    og.Text,
					Options: goptions,
				})
			}

			actions = append(actions, slack.AttachmentAction{
				Name:            ac.Name,
				Text:            ac.Text,
				Style:           ac.Style,
				Type:            ac.Type,
				Value:           ac.Value,
				DataSource:      ac.DataSource,
				MinQueryLength:  ac.MinQueryLength,
				Options:         options,
				SelectedOptions: selectedOptions,
				OptionGroups:    optionGroups,
				Confirm:         (*slack.ConfirmationField)(ac.Confirm),
				URL:             ac.URL,
			})
		}

		attachments = append(attachments, slack.Attachment{
			Color:      a.Color,
			Fallback:   a.Fallback,
			CallbackID: a.CallbackID,
			ID:         a.ID,
			AuthorID:   a.AuthorID,
			AuthorName: a.AuthorName,
			AuthorLink: a.AuthorLink,
			AuthorIcon: a.AuthorIcon,
			Title:      a.Title,
			TitleLink:  a.TitleLink,
			Pretext:    a.Pretext,
			Text:       a.Text,
			ImageURL:   a.ImageURL,
			ThumbURL:   a.ThumbURL,
			Fields:     fields,
			Actions:    actions,
			MarkdownIn: a.MarkdownIn,
			Footer:     a.Footer,
			FooterIcon: a.FooterIcon,
			Ts:         a.Ts,
		})
	}

	return &slack.WebhookMessage{
		Username:        m.Spec.Username,
		IconEmoji:       m.Spec.IconEmoji,
		IconURL:         m.Spec.IconURL,
		Channel:         m.Spec.Channel,
		ThreadTimestamp: m.Spec.ThreadTimestamp,
		Text:            m.Spec.Text,
		Attachments:     attachments,
	}
}

type MessageSpec struct {
	Username        string             `json:"username,omitempty"`
	IconEmoji       string             `json:"iconEmoji,omitempty"`
	IconURL         string             `json:"iconUrl,omitempty"`
	Channel         string             `json:"channel,omitempty"`
	ThreadTimestamp string             `json:"threadTs,omitempty"`
	Text            string             `json:"text,omitempty"`
	Attachments     []slack.Attachment `json:"attachments,omitempty"`
}

func (m *MessageSpec) Validate() error {
	if m.Channel == "" {
		return errors.New("channel must be specified")
	}

	if m.Text == "" && len(m.Attachments) == 0 {
		return errors.New("text or attachments must be specified")
	}

	return nil
}

type Attachment struct {
	Color    string `json:"color,omitempty"`
	Fallback string `json:"fallback"`

	CallbackID string `json:"callbackID,omitempty"`
	ID         int    `json:"id,omitempty"`

	AuthorID      string `json:"authorID,omitempty"`
	AuthorName    string `json:"authorName,omitempty"`
	AuthorSubname string `json:"authorSubname,omitempty"`
	AuthorLink    string `json:"authorLink,omitempty"`
	AuthorIcon    string `json:"authorIcon,omitempty"`

	Title     string `json:"title,omitempty"`
	TitleLink string `json:"titleLink,omitempty"`
	Pretext   string `json:"pretext,omitempty"`
	Text      string `json:"text"`

	ImageURL string `json:"imageURL,omitempty"`
	ThumbURL string `json:"thumbURL,omitempty"`

	Fields     []AttachmentField  `json:"fields,omitempty"`
	Actions    []AttachmentAction `json:"actions,omitempty"`
	MarkdownIn []string           `json:"mrkdwnIn,omitempty"`

	Footer     string `json:"footer,omitempty"`
	FooterIcon string `json:"footerIcon,omitempty"`

	Ts json.Number `json:"ts,omitempty"`
}

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type AttachmentAction struct {
	Name            string                        `json:"name"`
	Text            string                        `json:"text"`
	Style           string                        `json:"style,omitempty"`
	Type            string                        `json:"type"`
	Value           string                        `json:"value,omitempty"`
	DataSource      string                        `json:"dataSource,omitempty"`
	MinQueryLength  int                           `json:"minQueryLength,omitempty"`
	Options         []AttachmentActionOption      `json:"options,omitempty"`
	SelectedOptions []AttachmentActionOption      `json:"selectedOptions,omitempty"`
	OptionGroups    []AttachmentActionOptionGroup `json:"optionGroups,omitempty"`
	Confirm         *ConfirmationField            `json:"confirm,omitempty"`
	URL             string                        `json:"url,omitempty"`
}

type AttachmentActionOption struct {
	Text        string `json:"text"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

type AttachmentActionOptionGroup struct {
	Text    string                   `json:"text"`
	Options []AttachmentActionOption `json:"options"`
}

type ConfirmationField struct {
	Title       string `json:"title,omitempty"`
	Text        string `json:"text"`
	OkText      string `json:"okText,omitempty"`
	DismissText string `json:"dismissText,omitempty"`
}

type MessageStatus struct {
	SendTime int64 `json:"sendTime"`
}
