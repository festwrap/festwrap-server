package errors

type CannotRetrieveSetlistError struct {
	message string
}

func NewCannotRetrieveSetlistError(message string) error {
	return &CannotRetrieveSetlistError{message: message}
}

func (e *CannotRetrieveSetlistError) Error() string {
	return e.message
}
