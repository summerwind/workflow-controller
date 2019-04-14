package webhook

type AdmissionRequest struct {
	UID string `json:"uid"`
}

type AdmissionResponse struct {
	UID     string `json:"uid"`
	Allowed bool   `json:"allowed"`
}
