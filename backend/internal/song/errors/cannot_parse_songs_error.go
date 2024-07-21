package errors

type CannotParseSongsError struct {
	message string
}

func (e *CannotParseSongsError) Error() string {
	return e.message
}

func NewCannotParseSongsError(message string) error {
	return &CannotParseSongsError{message: message}
}
