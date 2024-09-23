package errors

type CannotAddSongsToPlaylistError struct {
	message string
}

func NewCannotAddSongsToPlaylistError(message string) error {
	return &CannotAddSongsToPlaylistError{message: message}
}

func (e *CannotAddSongsToPlaylistError) Error() string {
	return e.message
}
