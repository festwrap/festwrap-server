package serialization

import (
	"io"
)

type FakeEncoder[T any] struct {
	err error
}

func (e *FakeEncoder[T]) SetError(err error) {
	e.err = err
}

func (e FakeEncoder[T]) Encode(w io.Writer, object T) error {
	return e.err
}
