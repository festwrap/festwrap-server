package serialization

import (
	"festwrap/internal/testtools"
	"testing"
)

func TestBaseSerializerReturnsExpectedOutput(t *testing.T) {
	serializer := NewJsonSerializer[Object]()

	object := serializableObject()
	actual, err := serializer.Serialize(object)

	expected := serializableObjectBytes()
	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, actual, expected)
}

func TestBaseSerializerReturnsErrorOnNonSerializableObject(t *testing.T) {
	serializer := NewJsonSerializer[NonSerializableObject]()

	object := nonSerializableObject()
	_, err := serializer.Serialize(object)

	testtools.AssertErrorIsNotNil(t, err)
}
