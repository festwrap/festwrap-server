package serialization

type Serializer[T any] interface {
	Serialize(input T) ([]byte, error)
}
