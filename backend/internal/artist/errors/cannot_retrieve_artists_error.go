package errors

type CannotRetrieveArtistsError struct {
	message string
}

func (e *CannotRetrieveArtistsError) Error() string {
	return e.message
}

func NewCannotRetrieveArtistsError(message string) error {
	return &CannotRetrieveArtistsError{message: message}
}
