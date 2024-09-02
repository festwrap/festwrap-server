package serialization

type FakeDeserializer[T any] struct {
	args     []byte
	response *T
	err      error
}

func (s *FakeDeserializer[T]) GetArgs() []byte {
	return s.args
}

func (s *FakeDeserializer[T]) SetResponse(response *T) {
	s.response = response
}

func (s *FakeDeserializer[T]) SetError(err error) {
	s.err = err
}

func (s *FakeDeserializer[T]) Deserialize(bytes []byte) (*T, error) {
	s.args = bytes
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}
