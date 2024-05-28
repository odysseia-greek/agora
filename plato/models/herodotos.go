package models

type Section struct {
	Key string `json:"key"`
}

type Reference struct {
	Key      string    `json:"key"`
	Sections []Section `json:"sections"`
}

type ESBook struct {
	Key        string      `json:"key"`
	References []Reference `json:"references"`
}

type ESAuthor struct {
	Key   string   `json:"key"`
	Books []ESBook `json:"books"`
}

type AggregationResult struct {
	Authors []ESAuthor `json:"authors"`
}

type Rhema struct {
	Greek        string   `json:"greek"`
	Translations []string `json:"translations"`
	Section      string   `json:"section"`
}
type Text struct {
	Author          string  `json:"author"`
	Book            string  `json:"book"`
	Type            string  `json:"type"`
	Reference       string  `json:"reference"`
	PerseusTextLink string  `json:"perseusTextLink"`
	Rhemai          []Rhema `json:"rhemai"`
}

type AnalyzeTextRequest struct {
	// example: Ἀθηναῖος
	// required: true
	Rootword string `json:"rootword"`
}

type AnalyzeTextResponse struct {
	// example: Ἀθηναῖος
	// required: true
	Rootword string `json:"rootword"`
	// example: ["Ἀθηναῖος"]
	// required: true
	Conjugations []string        `json:"conjugations"`
	Results      []AnalyzeResult `json:"texts"`
}

type AnalyzeResult struct {
	// example: text/author=herodotos&book=histories&reference=1.1
	ReferenceLink string `json:"referenceLink"`
	Text          Rhema  `json:"text"`
}

type CreateTextRequest struct {
	// example: Herodotos
	// required: true
	Author string `json:"author"`
	// example: Histories
	// required: true
	Book string `json:"book"`
	// example: 1.1
	// required: true
	Reference string `json:"reference"`
	// example: a
	// required: false
	Section string `json:"section"`
}

type CheckTextRequest struct {
	Translations []Translations `json:"translations"`
	// example: Herodotos
	// required: true
	Author string `json:"author"`
	// example: Histories
	// required: true
	Book string `json:"book"`
	// example: 1.1
	// required: true
	Reference string `json:"reference"`
}

type Translations struct {
	// example: a
	// required: false
	Section string `json:"section"`
	// example: this is an example sentence
	// required: true
	Translation string `json:"translation"`
}

type CheckTextResponse struct {
	AverageLevenshteinPercentage string          `json:"averageLevenshteinPercentage"`
	Sections                     []AnswerSection `json:"sections"`
	PossibleTypos                []Typo          `json:"possibleTypos"`
}

type Typo struct {
	Source   string `json:"source"`
	Provided string `json:"provided"`
}

type AnswerSection struct {
	// example: a
	// required: false
	Section string `json:"section"`
	// example: 9.09
	// required: true
	LevenshteinPercentage string `json:"levenshteinPercentage"`
	// example: Such a step would not be condemned either by the gods who received our oaths,
	// required: true
	QuizSentence string `json:"quizSentence"`
	// example: this is an example answer"
	// required: true
	AnswerSentence string `json:"answerSentence"`
}
