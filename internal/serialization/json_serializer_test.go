package serialization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseSerializerReturnsExpectedOutput(t *testing.T) {
	serializer := NewJsonSerializer[Object]()

	object := serializableObject()
	actual, err := serializer.Serialize(object)

	expected := serializableObjectBytes()
	assert.Nil(t, err)
	assert.Equal(t, actual, expected)
}

func TestBaseSerializerReturnsErrorOnNonSerializableObject(t *testing.T) {
	serializer := NewJsonSerializer[NonSerializableObject]()

	object := nonSerializableObject()
	_, err := serializer.Serialize(object)

	assert.NotNil(t, err)
}
