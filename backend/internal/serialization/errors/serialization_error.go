package errors

type SerializationError struct {
	message string
}

func NewSerializationError(message string) error {
	return &SerializationError{message: message}
}

func (e *SerializationError) Error() string {
	return e.message
}
