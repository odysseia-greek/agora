package models

type EdgecaseRequest struct {
	// example: ferw
	// required: true
	Rootword string `json:"rootword"`
}

type EdgecaseResponse struct {
	OriginalWord   string  `json:"originalWord"`
	GreekWord      string  `json:"greekWord"`
	StrongPassword string  `json:"strongPassword"`
	SimilarWords   []Meros `json:"similarWords,omitempty"`
}
