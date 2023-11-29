package models

// swagger:model
// ErrorModel is the base model used for handling errors
type ErrorModel struct {
	// required: true
	// example: 94374b4f-3dda-4ffb-b33b-2cb6ba092b84
	UniqueCode string `json:"uniqueCode"`
}

// swagger:model
// ValidationMessages messages used in validation error
type ValidationMessages struct {
	// example: word
	Field string `json:"validationField"`
	// example: cannot be empty
	Message string `json:"validationMessage"`
}

// swagger:model
// ValidationError validation errors occur when data is malformed
type ValidationError struct {
	ErrorModel
	Messages []ValidationMessages `json:"errorModel"`
}

// swagger:model
type NotFoundError struct {
	ErrorModel
	Message NotFoundMessage `json:"errorModel"`
}

func (m *NotFoundError) Error() string {
	return m.Error()
}

// swagger:model
type NotFoundMessage struct {
	// example: query for obscura
	Type string `json:"type"`
	// example: produced 0 results
	Reason string `json:"reason"`
}

// swagger:model
type ElasticSearchError struct {
	ErrorModel
	Message ElasticErrorMessage `json:"errorModel"`
}

func (m *ElasticSearchError) Error() string {
	return m.Error()
}

// swagger:model
type ElasticErrorMessage struct {
	ElasticError string `json:"elasticError"`
}

// MethodMessages messages used in method error
// swagger:model
type MethodMessages struct {
	// example: GET
	Methods string `json:"allowedMethods"`
	// example: Method DELETE not allowed at this endpoint
	Message string `json:"methodError"`
}

// swagger:model
type MethodError struct {
	ErrorModel
	Messages []MethodMessages `json:"errorModel"`
}
