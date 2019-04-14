package message

import (
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

type MessageSpec struct {
	Channel     string             `json:"channel"`
	Text        string             `json:"text,omitempty"`
	Username    string             `json:"username,omitempty"`
	Markdown    bool               `json:"mrkdwn"`
	Attachments []slack.Attachment `json:"attachments,omitempty"`
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

type MessageStatus struct {
	SendTime int64 `json:"sendTime"`
}

type State struct {
	Resource *Message `json:"resource"`
}
