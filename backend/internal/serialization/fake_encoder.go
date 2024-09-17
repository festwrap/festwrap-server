package serialization

import (
	"io"
)

type EncodeArgs[T any] struct {
	Writer io.Writer
	Object T
}

type FakeEncoder[T any] struct {
	encodeArgs EncodeArgs[T]
	err        error
}

func (e FakeEncoder[T]) GetEncodeArgs() EncodeArgs[T] {
	return e.encodeArgs
}

func (e *FakeEncoder[T]) SetError(err error) {
	e.err = err
}

func (e FakeEncoder[T]) Encode(w io.Writer, object T) error {
	return e.err
}
