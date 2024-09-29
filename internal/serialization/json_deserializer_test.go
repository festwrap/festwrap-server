package serialization

import (
	"festwrap/internal/testtools"
	"testing"
)

func TestJsonDeserializerProducesExpectedResult(t *testing.T) {
	deserializer := NewJsonDeserializer[Object]()

	actual, err := deserializer.Deserialize(serializableObjectBytes())

	expected := serializableObject()
	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, *actual, expected)
}

func TestJsonDeserializerReturnsErrorOnNonJsonInput(t *testing.T) {
	deserializer := NewJsonDeserializer[Object]()

	nonJsonBytes := []byte(`"something": "something"`)
	_, err := deserializer.Deserialize(nonJsonBytes)

	testtools.AssertErrorIsNotNil(t, err)
}
