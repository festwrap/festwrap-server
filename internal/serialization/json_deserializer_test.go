package serialization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonDeserializerProducesExpectedResult(t *testing.T) {
	deserializer := NewJsonDeserializer[Object]()

	var actual Object
	err := deserializer.Deserialize(serializableObjectBytes(), &actual)

	expected := serializableObject()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestJsonDeserializerReturnsErrorOnNonJsonInput(t *testing.T) {
	deserializer := NewJsonDeserializer[Object]()

	var object Object
	nonJsonBytes := []byte(`"something": "something"`)
	err := deserializer.Deserialize(nonJsonBytes, &object)

	assert.NotNil(t, err)
}
