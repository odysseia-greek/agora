package models

import "encoding/json"

func UnmarshalRhema(data []byte) (RhemaSource, error) {
	var r RhemaSource
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RhemaSource) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type RhemaSource struct {
	// example: Herodotos
	// required: true
	Author string `json:"author"`
	// example: ὡς δέ οἱ ταῦτα ἔδοξε, καὶ ἐποίεε κατὰ τάχος·
	// required: true
	Greek string `json:"greek"`
	// example: ["first translation", "second translation"]
	// required: true
	Translations []string `json:"translations"`
	// example: 1
	// required: true
	Book int64 `json:"book"`
	// example: 1
	// required: true
	Chapter int64 `json:"chapter"`
	// example: 1
	// required: true
	Section int64 `json:"section"`
	// example: https://externallink
	// required: true
	PerseusTextLink string `json:"perseusTextLink"`
}

// swagger:model
type Rhema struct {
	Rhemai []RhemaSource `json:"rhemai"`
}

func UnmarshalAuthors(data []byte) (Author, error) {
	var r Author
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Author) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type Authors struct {
	Authors []Author `json:"authors"`
}

// swagger:model
type Author struct {
	// example: herodotos
	// required: true
	Author string `json:"author"`
}

// swagger:model
type Books struct {
	Books []Book `json:"books"`
}

// swagger:model
type Book struct {
	// example: 2
	Book int64 `json:"book"`
}

// swagger:model
type CreateSentenceResponse struct {
	// example: ὡς δέ οἱ ταῦτα ἔδοξε, καὶ ἐποίεε κατὰ τάχος·
	// required: true
	Sentence string `json:"sentence"`
	// example: fd4TlogBC__qOhD2dK31
	// required: true
	SentenceId string `json:"sentenceId"`
}

func (r *CheckSentenceRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type CheckSentenceRequest struct {
	// example: fd4TlogBC__qOhD2dK31
	// required: true
	SentenceId string `json:"sentenceId"`
	// example: this is an example sentence
	// required: true
	ProvidedSentence string `json:"answerSentence"`
	// example: herodotos
	// required: true
	Author string `json:"author"`
}

// swagger:model
type CheckSentenceResponse struct {
	// example: 9.09
	// required: true
	LevenshteinPercentage string `json:"levenshteinPercentage"`
	// example: Such a step would not be condemned either by the gods who received our oaths,
	// required: true
	QuizSentence string `json:"quizSentence"`
	// example: this is an example answer"
	// required: true
	AnswerSentence string `json:"answerSentence"`
	// example: ["Such", "condemned"]
	// required: true
	SplitQuizSentence []string `json:"splitQuizSentence"`
	// example: ["this", "example"]
	// required: true
	SplitAnswerSentence []string          `json:"splitAnswerSentence"`
	MatchingWords       []MatchingWord    `json:"matchingWords,omitempty"`
	NonMatchingWords    []NonMatchingWord `json:"nonMatchingWords,omitempty"`
}

// swagger:model
type MatchingWord struct {
	// example: thiswordisinthetext
	// required: true
	Word string `json:"word"`
	// example: 4
	// required: true
	SourceIndex int `json:"sourceIndex"`
}

// swagger:model
type NonMatchingWord struct {
	// example: step
	// required: true
	Word string `json:"word"`
	// example: 3
	// required: true
	SourceIndex int     `json:"sourceIndex"`
	Matches     []Match `json:"matches"`
}

// swagger:model
type Match struct {
	// example: superduperword
	// required: true
	Match string `json:"match"`
	// example: 4
	// required: true
	Levenshtein int `json:"levenshtein"`
	// example: 3
	// required: true
	AnswerIndex int `json:"answerIndex"`
	// example: 25.00
	// required: true
	Percentage string `json:"percentage"`
}
