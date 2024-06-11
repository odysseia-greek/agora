package models

const (
	MEDIA       string = "media"
	AUTHORBASED string = "authorbased"
	DIALOGUE    string = "dialogue"
)

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
	Theme    string         `json:"theme,omitempty"`
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

type AggregateResult struct {
	Aggregates []Aggregate `json:"aggregates"`
}

type Aggregate struct {
	HighestSet string `json:"highestSet"`
	Name       string `json:"name"`
}

type CreationRequest struct {
	Theme    string `json:"theme"`
	Set      string `json:"set"`
	QuizType string `json:"quizType"`
}

type AnswerRequest struct {
	Theme         string            `json:"theme"`
	Set           string            `json:"set"`
	QuizType      string            `json:"quizType"`
	Comprehensive bool              `json:"comprehensive,omitempty"`
	Answer        string            `json:"answer"`
	Dialogue      []DialogueContent `json:"dialogue,omitempty"`
	QuizWord      string            `json:"quizWord"`
}

type QuizResponse struct {
	QuizItem string    `json:"quizItem"`
	Options  []Options `json:"options,omitempty"`
}

type Options struct {
	Option   string `json:"quizWord"`
	AudioUrl string `json:"audioUrl,omitempty"`
	ImageUrl string `json:"imageUrl,omitempty"`
}

type ComprehensiveResponse struct {
	Correct      bool                `json:"correct"`
	QuizWord     string              `json:"quizWord"`
	FoundInText  AnalyzeTextResponse `json:"foundInText,omitempty"`
	SimilarWords []Meros             `json:"similarWords,omitempty"`
	Progress     struct {
		TimesCorrect    int     `json:"timesCorrect"`
		TimesIncorrect  int     `json:"timesIncorrect"`
		AverageAccuracy float64 `json:"averageAccuracy"`
	} `json:"progress,omitempty"`
}

type DialogueAnswer struct {
	Percentage   float64              `json:"percentage"`
	Input        []DialogueContent    `json:"input"`
	Answer       []DialogueContent    `json:"answer"`
	InWrongPlace []DialogueCorrection `json:"wronglyPlaced"`
}

type DialogueCorrection struct {
	Translation  string `json:"translation"`
	Greek        string `json:"greek,omitempty"`
	Place        int    `json:"place,omitempty"`
	Speaker      string `json:"speaker,omitempty"`
	CorrectPlace int    `json:"correctPlace,omitempty"`
}

type QuizAttempt struct {
	Correct bool
	Set     string
	Theme   string
}
