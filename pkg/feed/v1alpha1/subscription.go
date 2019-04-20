package v1alpha1

import (
	"errors"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const timeMergin = 30

type Subscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubscriptionSpec   `json:"spec,omitempty"`
	Status SubscriptionStatus `json:"status,omitempty"`
}

func (s *Subscription) Validate() error {
	return s.Spec.Validate()
}

func (s *Subscription) NeedsUpdate() bool {
	if s.Status.LastFetchTime+timeMergin > time.Now().Unix() {
		return false
	}

	return true
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
	LastFetchTime int64 `json:"lastFetchTime"`
}
