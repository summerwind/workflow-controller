package webhook

type AdmissionRequest struct {
	UID    string      `json:"uid"`
	Object interface{} `json:"object"`
}

type AdmissionResponse struct {
	UID     string `json:"uid"`
	Allowed bool   `json:"allowed"`
}
