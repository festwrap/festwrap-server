package errors

type CannotCreatePlaylistError struct {
	message string
}

func (e *CannotCreatePlaylistError) Error() string {
	return e.message
}

func NewCannotCreatePlaylistError(message string) error {
	return &CannotCreatePlaylistError{message: message}
}
