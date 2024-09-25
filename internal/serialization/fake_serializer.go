package serialization

type FakeSerializer[T any] struct {
	response []byte
	err      error
	args     T
}

func (s *FakeSerializer[T]) SetResponse(response []byte) {
	s.response = response
}

func (s *FakeSerializer[T]) SetError(err error) {
	s.err = err
}

func (s *FakeSerializer[T]) GetArgs() T {
	return s.args
}

func (s *FakeSerializer[T]) Serialize(input T) ([]byte, error) {
	s.args = input
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}
