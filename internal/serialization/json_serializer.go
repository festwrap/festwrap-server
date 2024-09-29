package serialization

import (
	"encoding/json"
)

type JsonSerializer[T any] struct{}

func (s *JsonSerializer[T]) Serialize(object T) ([]byte, error) {
	return json.Marshal(object)
}

func NewJsonSerializer[T any]() JsonSerializer[T] {
	return JsonSerializer[T]{}
}
