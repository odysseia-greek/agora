package models

import "encoding/json"

func UnmarshalSolonCreationRequest(data []byte) (SolonCreationRequest, error) {
	var r SolonCreationRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SolonCreationRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type SolonCreationRequest struct {
	// example: api
	// required: true
	Role string `json:"roles"`
	// example: [grammar dictionary]
	// required: true
	Access []string `json:"access"`
	// example: dionysios-544c584d7f-6sp6x
	// required: true
	PodName string `json:"podName"`
	// example: dionysios
	// required: true
	Username string `json:"username"`
}

// swagger:model
type SolonResponse struct {
	// example: true
	// required: true
	UserCreated bool `json:"userCreated"`
	// example: true
	// required: true
	SecretCreated bool `json:"secretCreated"`
}

func (r *SolonResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *TokenResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type TokenResponse struct {
	// example: s.1283745jdf83r3
	// required: true
	Token string `json:"token"`
}
