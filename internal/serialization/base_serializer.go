package serialization

import (
	"encoding/json"
)

type BaseSerializer[T any] struct{}

func (s *BaseSerializer[T]) Serialize(object T) ([]byte, error) {
	return json.Marshal(object)
}

func NewBaseSerializer[T any]() BaseSerializer[T] {
	return BaseSerializer[T]{}
}
