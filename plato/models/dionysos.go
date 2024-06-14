package models

type DeclensionConfig struct {
	Declensions []Declension
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

type FoundRules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Rule        string   `json:"rule,omitempty"`
	SearchTerms []string `json:"searchTerm,omitempty"`
}

type DeclensionTranslationResults struct {
	Results []Result `json:"results"`
}

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
	// example: [throw, to throw]
	// required: true
	Translation []string `json:"translations"`
}

func (r *DeclensionTranslationResults) RemoveIndex(index int) {
	r.Results = append(r.Results[:index], r.Results[index+1:]...)
}
