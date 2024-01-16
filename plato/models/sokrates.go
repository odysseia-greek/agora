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

type AuthorBasedQuiz struct {
	QuizMetadata struct {
		Language string `json:"language"`
	} `json:"quizMetadata"`
	QuizType string               `json:"quizType"`
	Theme    string               `json:"theme,omitempty"`
	Set      int                  `json:"set,omitempty"`
	Content  []AuthorBasedContent `json:"content"`
	Progress struct {
		TimesCorrect    int     `json:"timesCorrect"`
		TimesIncorrect  int     `json:"timesIncorrect"`
		AverageAccuracy float64 `json:"averageAccuracy"`
	} `json:"progress,omitempty"`
}

type AuthorBasedContent struct {
	Translation     string  `json:"translation"`
	TimesCorrect    int     `json:"timesCorrect,omitempty"`
	TimesIncorrect  int     `json:"timesIncorrect,omitempty"`
	AverageAccuracy float64 `json:"averageAccuracy,omitempty"`
	Greek           string  `json:"greek,omitempty"`
}

type MediaContent struct {
	Translation string `json:"translation"`
	Greek       string `json:"greek,omitempty"`
	ImageURL    string `json:"imageURL,omitempty"`
	AudioFile   string `json:"audioFile,omitempty"`
}

type DialogueContent struct {
	Translation string `json:"translation"`
	Greek       string `json:"greek,omitempty"`
	Place       int    `json:"place,omitempty"`
	Speaker     string `json:"speaker,omitempty"`
}

type MediaQuiz struct {
	QuizMetadata struct {
		Language string `json:"language"`
	} `json:"quizMetadata"`
	QuizType string         `json:"quizType"`
	Set      int            `json:"set,omitempty"`
	Content  []MediaContent `json:"content"`
}

type DialogueQuiz struct {
	QuizMetadata struct {
		Language string `json:"language"`
	} `json:"quizMetadata"`
	Theme    string            `json:"theme,omitempty"`
	QuizType string            `json:"quizType"`
	Set      int               `json:"set,omitempty"`
	Dialogue Dialogue          `json:"dialogue,omitempty"`
	Content  []DialogueContent `json:"content"`
}

type Dialogue struct {
	Introduction string `json:"introduction"`
	Speakers     []struct {
		Name        string `json:"name"`
		Shorthand   string `json:"shorthand"`
		Translation string `json:"translation"`
	} `json:"speakers"`
	Section       string `json:"section"`
	LinkToPerseus string `json:"linkToPerseus"`
}
