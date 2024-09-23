package errors

type CannotCreatePlaylistError struct {
	message string
}

func NewCannotCreatePlaylistError(message string) error {
	return &CannotCreatePlaylistError{message: message}
}

func (e *CannotCreatePlaylistError) Error() string {
	return e.message
}
