package validator

// Validator is an interface for validating a value
type Validator interface {
	Validate() error
}
