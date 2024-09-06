package errors

type DeserializationError struct {
	message string
}

func (e *DeserializationError) Error() string {
	return e.message
}

func NewDeserializationError(message string) error {
	return &DeserializationError{message: message}
}
