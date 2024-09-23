package errors

type CannotRetrieveArtistsError struct {
	message string
}

func NewCannotRetrieveArtistsError(message string) error {
	return &CannotRetrieveArtistsError{message: message}
}

func (e *CannotRetrieveArtistsError) Error() string {
	return e.message
}
