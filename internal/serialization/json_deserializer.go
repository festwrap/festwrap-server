package serialization

import (
	"encoding/json"
)

type JsonDeserializer[T any] struct{}

func (s JsonDeserializer[T]) Deserialize(bytes []byte) (*T, error) {
	var dest T
	err := json.Unmarshal(bytes, &dest)
	return &dest, err
}

func NewJsonDeserializer[T any]() JsonDeserializer[T] {
	return JsonDeserializer[T]{}
}
