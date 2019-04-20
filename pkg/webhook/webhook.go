package webhook

import (
	"encoding/json"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AdmissionRequest struct {
	v1beta1.AdmissionRequest
}

func NewAdmissionRequest() *AdmissionRequest {
	return &AdmissionRequest{
		AdmissionRequest: v1beta1.AdmissionRequest{},
	}
}

func (r *AdmissionRequest) GetObject(object interface{}) error {
	return json.Unmarshal(r.Object.Raw, object)
}

type AdmissionResponse struct {
	v1beta1.AdmissionResponse
}

func NewAdmissionResponse(req *AdmissionRequest) *AdmissionResponse {
	return &AdmissionResponse{
		v1beta1.AdmissionResponse{
			UID: req.UID,
		},
	}
}

func (r *AdmissionResponse) SetSuccess() {
	r.Allowed = true
}

func (r *AdmissionResponse) SetFailure(msg string) {
	r.Allowed = false
	r.Result = &metav1.Status{
		Status: "Failure",
		Reason: metav1.StatusReason(msg),
	}
}
