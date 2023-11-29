package models

import "encoding/json"

type DeclensionConfig struct {
	Declensions []Declension
}

func UnmarshalDeclension(data []byte) (Declension, error) {
	var r Declension
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Declension) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Declension struct {
	Name        string              `json:"name"`
	Type        string              `json:"type,omitempty"`
	Dialect     string              `json:"dialect"`
	Declensions []DeclensionElement `json:"declensions"`
}

type DeclensionElement struct {
	Declension string   `json:"declension"`
	RuleName   string   `json:"ruleName"`
	SearchTerm []string `json:"searchTerm"`
}

func UnmarshalFoundRules(data []byte) (FoundRules, error) {
	var r FoundRules
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FoundRules) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FoundRules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Rule        string   `json:"rule,omitempty"`
	SearchTerms []string `json:"searchTerm,omitempty"`
}

func UnmarshalDeclensionTranslationResults(data []byte) (DeclensionTranslationResults, error) {
	var r DeclensionTranslationResults
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DeclensionTranslationResults) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type DeclensionTranslationResults struct {
	Results []Result `json:"results"`
}

// swagger:model
type Result struct {
	// example: ἔβαλλε
	// required: true
	Word string `json:"word"`
	// example: 3th sing - impf - ind - act
	// required: true
	Rule string `json:"rule"`
	// example: βαλλω
	// required: true
	RootWord string `json:"rootWord"`
	// example: throw
	// required: true
	Translation string `json:"translation"`
}

func (r *DeclensionTranslationResults) RemoveIndex(index int) {
	r.Results = append(r.Results[:index], r.Results[index+1:]...)
}
