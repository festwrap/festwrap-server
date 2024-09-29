package serialization

import (
	"encoding/json"
)

type JsonDeserializer[T any] struct{}

func (s JsonDeserializer[T]) Deserialize(bytes []byte, dest *T) error {
	return json.Unmarshal(bytes, dest)
}

func NewJsonDeserializer[T any]() JsonDeserializer[T] {
	return JsonDeserializer[T]{}
}
