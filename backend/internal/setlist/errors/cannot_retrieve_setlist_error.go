package errors

type CannotRetrieveSetlistError struct {
	message string
}

func (e *CannotRetrieveSetlistError) Error() string {
	return e.message
}

func NewCannotRetrieveSetlistError(message string) error {
	return &CannotRetrieveSetlistError{message: message}
}
