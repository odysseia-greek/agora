package models

import "encoding/json"

func UnmarshalWord(data []byte) (Word, error) {
	var r Word
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Word) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type Word struct {
	// example: aristophanes
	// required: true
	Method string `json:"method"`
	// example: frogs
	// required: true
	Category string `json:"category"`
	// example: ὄνος
	// required: true
	Greek string `json:"greek"`
	// example: donkey
	// required: true
	Translation string `json:"translation"`
	// example: 1
	// required: true
	Chapter int64 `json:"chapter"`
}

func UnmarshalLogos(data []byte) (Logos, error) {
	var r Logos
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Logos) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type Logos struct {
	Logos []Word `json:"logos"`
}
