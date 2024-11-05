package models

const (
	MEDIA       string = "media"
	AUTHORBASED string = "authorbased"
	MULTICHOICE string = "multiplechoice"
	DIALOGUE    string = "dialogue"
)

type AuthorbasedQuiz struct {
	QuizType               string                 `json:"quizType"`
	Theme                  string                 `json:"theme"`
	Set                    int                    `json:"set"`
	Segment                string                 `json:"segment,omitempty"`
	Reference              string                 `json:"reference"`
	FullSentence           string                 `json:"fullSentence"`
	Translation            string                 `json:"translation"`
	GrammarQuestionOptions GrammarQuestionOptions `json:"grammarQuestionOptions"`
	Content                []AuthorBasedContent   `json:"content"`
}

type AuthorBasedContent struct {
	Greek               string            `json:"greek"`
	Translation         string            `json:"translation"`
	WordsInText         []string          `json:"wordsInText"`
	HasGrammarQuestions bool              `json:"hasGrammarQuestions"`
	GrammarQuestions    []GrammarQuestion `json:"grammarQuestions,omitempty"`
}

type GrammarQuestion struct {
	CorrectAnswer    string `json:"correctAnswer"`
	TypeOfWord       string `json:"typeOfWord"`
	WordInText       string `json:"wordInText"`
	ExtraInformation string `json:"extraInformation"`
}

type GrammarQuizAdded struct {
	CorrectAnswer    string    `json:"correctAnswer"`
	WordInText       string    `json:"wordInText"`
	ExtraInformation string    `json:"extraInformation"`
	Options          []Options `json:"options,omitempty"`
}

type GrammarQuestionOptions struct {
	Nouns []string `json:"nouns"`
	Verbs []string `json:"verbs"`
	Misc  []string `json:"misc"`
}

type MultipleChoiceQuiz struct {
	QuizMetadata struct {
		Language string `json:"language"`
	} `json:"quizMetadata"`
	QuizType string                  `json:"quizType"`
	Theme    string                  `json:"theme,omitempty"`
	Set      int                     `json:"set,omitempty"`
	Content  []MultipleChoiceContent `json:"content"`
	Progress struct {
		TimesCorrect    int     `json:"timesCorrect"`
		TimesIncorrect  int     `json:"timesIncorrect"`
		AverageAccuracy float64 `json:"averageAccuracy"`
	} `json:"progress,omitempty"`
}

type MultipleChoiceContent struct {
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
	Segment  string         `json:"segment,omitempty"`
	Content  []MediaContent `json:"content"`
	Progress struct {
		TimesCorrect    int     `json:"timesCorrect"`
		TimesIncorrect  int     `json:"timesIncorrect"`
		AverageAccuracy float64 `json:"averageAccuracy"`
	} `json:"progress,omitempty"`
}

type DialogueQuiz struct {
	QuizMetadata struct {
		Language string `json:"language"`
	} `json:"quizMetadata"`
	Theme     string            `json:"theme,omitempty"`
	QuizType  string            `json:"quizType"`
	Set       int               `json:"set,omitempty"`
	Segment   string            `json:"segment,omitempty"`
	Reference string            `json:"reference,omitempty"`
	Dialogue  Dialogue          `json:"dialogue,omitempty"`
	Content   []DialogueContent `json:"content"`
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

type Segment struct {
	Name   string  `json:"name"`
	MaxSet float64 `json:"maxSet"`
}

type Theme struct {
	Name     string    `json:"name"`
	Segments []Segment `json:"segments"`
}

type AggregatedOptions struct {
	Themes []Theme `json:"themes"`
}

type CreationRequest struct {
	Theme        string   `json:"theme"`
	Set          string   `json:"set"`
	Segment      string   `json:"segment,omitempty"`
	QuizType     string   `json:"quizType"`
	Order        string   `json:"order"`
	ExcludeWords []string `json:"excludeWords"`
}

type AnswerRequest struct {
	Theme         string            `json:"theme"`
	Set           string            `json:"set"`
	QuizType      string            `json:"quizType"`
	Segment       string            `json:"segment,omitempty"`
	Comprehensive bool              `json:"comprehensive,omitempty"`
	Answer        string            `json:"answer"`
	Dialogue      []DialogueContent `json:"dialogue,omitempty"`
	QuizWord      string            `json:"quizWord"`
}

type AuthorbasedQuizResponse struct {
	FullSentence string             `json:"fullSentence"`
	Translation  string             `json:"translation"`
	Reference    string             `json:"reference"`
	Quiz         QuizResponse       `json:"quiz"`
	GrammarQuiz  []GrammarQuizAdded `json:"grammarQuiz,omitempty"`
}

type QuizResponse struct {
	QuizItem      string    `json:"quizItem"`
	NumberOfItems int       `json:"numberOfItems"`
	Options       []Options `json:"options,omitempty"`
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

type AuthorBasedResponse struct {
	Correct     bool     `json:"correct"`
	QuizWord    string   `json:"quizWord"`
	WordsInText []string `json:"wordsInText,omitempty"`
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
	Correct  bool
	Set      string
	Theme    string
	QuizType string
}
