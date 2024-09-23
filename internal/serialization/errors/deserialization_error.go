package errors

type DeserializationError struct {
	message string
}

func NewDeserializationError(message string) error {
	return &DeserializationError{message: message}
}

func (e *DeserializationError) Error() string {
	return e.message
}
