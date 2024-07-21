package errors

type CannotParseSetlistError struct {
	message string
}

func (e *CannotParseSetlistError) Error() string {
	return e.message
}

func NewCannotParseSetlistError(message string) error {
	return &CannotParseSetlistError{message: message}
}
