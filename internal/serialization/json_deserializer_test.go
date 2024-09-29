package serialization

import (
	"festwrap/internal/testtools"
	"testing"
)

func TestJsonDeserializerProducesExpectedResult(t *testing.T) {
	deserializer := NewJsonDeserializer[Object]()

	var actual Object
	err := deserializer.Deserialize(serializableObjectBytes(), &actual)

	expected := serializableObject()
	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, actual, expected)
}

func TestJsonDeserializerReturnsErrorOnNonJsonInput(t *testing.T) {
	deserializer := NewJsonDeserializer[Object]()

	var object Object
	nonJsonBytes := []byte(`"something": "something"`)
	err := deserializer.Deserialize(nonJsonBytes, &object)

	testtools.AssertErrorIsNotNil(t, err)
}
