package subscription

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Subscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubscriptionSpec   `json:"spec,omitempty"`
	Status SubscriptionStatus `json:"status,omitempty"`
}

func (s *Subscription) Validate() error {
	return s.Spec.Validate()
}

type SubscriptionSpec struct {
	URL string `json:"url"`
}

func (s *SubscriptionSpec) Validate() error {
	if s.URL == "" {
		return errors.New("url must be specified")
	}

	return nil
}

type SubscriptionStatus struct {
	LastFetchedTime int64 `json:"lastFetchedTime"`
}
