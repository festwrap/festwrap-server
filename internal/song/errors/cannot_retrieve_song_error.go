package errors

type CannotRetrieveSongError struct {
	message string
}

func NewCannotRetrieveSongError(message string) error {
	return &CannotRetrieveSongError{message: message}
}

func (e *CannotRetrieveSongError) Error() string {
	return e.message
}
