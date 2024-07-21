package errors

type CannotRetrieveSongError struct {
	message string
}

func (e *CannotRetrieveSongError) Error() string {
	return e.message
}

func NewCannotRetrieveSongError(message string) error {
	return &CannotRetrieveSongError{message: message}
}
