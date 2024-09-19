package serialization

import (
	"encoding/json"
	"io"
)

type Encoder[T any] interface {
	Encode(w io.Writer, object T) error
}

type JsonEncoder[T any] struct{}

func NewJsonEncoder[T any]() JsonEncoder[T] {
	return JsonEncoder[T]{}
}

func (e JsonEncoder[T]) Encode(w io.Writer, object T) error {
	return json.NewEncoder(w).Encode(object)
}
