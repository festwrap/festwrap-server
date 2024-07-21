package errors

type CannotAddSongsToPlaylistError struct {
	message string
}

func (e *CannotAddSongsToPlaylistError) Error() string {
	return e.message
}

func NewCannotAddSongsToPlaylistError(message string) error {
	return &CannotAddSongsToPlaylistError{message: message}
}
