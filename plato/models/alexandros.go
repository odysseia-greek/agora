package models

import "encoding/json"

func (r *Biblos) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Meros) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalBiblos(data []byte) (Biblos, error) {
	var r Biblos
	err := json.Unmarshal(data, &r)
	return r, err
}

func UnmarshalMeros(data []byte) (Meros, error) {
	var r Meros
	err := json.Unmarshal(data, &r)
	return r, err
}

type Biblos struct {
	Biblos []Meros `json:"biblos"`
}

type ExtendedResponse struct {
	Hits []Hit `json:"hits,omitempty"`
}

type Hit struct {
	Hit         Meros  `json:"hit"`
	FoundInText *Rhema `json:"foundInText,omitempty"`
}

type Meros struct {
	// example: ὄνος
	// required: true
	Greek string `json:"greek"`
	// example: an ass
	// required: true
	English string `json:"english"`
	// example: ezel
	// required: false
	Dutch string `json:"dutch,omitempty"`
	// required: false
	LinkedWord string `json:"linkedWord,omitempty"`
	// required: false
	Original string `json:"original,omitempty"`
}
