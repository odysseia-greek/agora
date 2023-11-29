package models

import "encoding/json"

// swagger:model
type Methods struct {
	Method []Method `json:"methods"`
}

// swagger:model
type Method struct {
	// example: aristophanes
	// required: true
	Method string `json:"method"`
}

// swagger:model
type Categories struct {
	Category []Category `json:"categories"`
}

// swagger:model
type Category struct {
	// example: frogs
	// required: true
	Category string `json:"category"`
}

// swagger:model
type LastChapterResponse struct {
	// example: 119
	// required: true
	LastChapter int64 `json:"lastChapter"`
}

// swagger:model
type QuizResponse struct {
	// example: ὄνος
	// required: true
	Question string `json:"question"`
	// example: donkey
	// required: true
	Answer string `json:"answer"`
	// example: ["donkey", "anotheranswer"]
	// required: true
	QuizQuestions []string `json:"quiz"`
}

func (r *CheckAnswerRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// swagger:model
type CheckAnswerRequest struct {
	// example: ὄνος
	// required: true
	QuizWord string `json:"quizWord"`
	// example: horse
	// required: true
	AnswerProvided string `json:"answerProvided"`
}

// swagger:model
type CheckAnswerResponse struct {
	// example: false
	// required: true
	Correct bool `json:"correct"`
	// example: ὄνος
	// required: true
	QuizWord      string `json:"quizWord"`
	Possibilities []Word `json:"possibilities"`
}
