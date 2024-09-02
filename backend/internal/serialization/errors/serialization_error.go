package errors

type SerializationError struct {
	message string
}

func (e *SerializationError) Error() string {
	return e.message
}

func NewSerializationError(message string) error {
	return &SerializationError{message: message}
}
