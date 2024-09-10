package errors

type ImageNotFoundError struct {
	message string
}

func (e *ImageNotFoundError) Error() string {
	return e.message
}

func NewImageNotFoundError(message string) error {
	return &ImageNotFoundError{message: message}
}
