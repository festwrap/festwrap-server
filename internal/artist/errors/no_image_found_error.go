package errors

type ImageNotFoundError struct {
	message string
}

func NewImageNotFoundError(message string) error {
	return &ImageNotFoundError{message: message}
}

func (e *ImageNotFoundError) Error() string {
	return e.message
}
