package serialization

import (
	"festwrap/internal/testtools"
	"testing"
)

type Object struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type NonSerializableObject struct {
	Function func() string
}

func TestBaseSerializerReturnsExpectedOutput(t *testing.T) {
	serializer := NewBaseSerializer[Object]()

	object := Object{Name: "myname", Value: 10}
	actual, err := serializer.Serialize(object)

	expected := []byte(`{"name":"myname","value":10}`)
	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, actual, expected)
}

func TestBaseSerializerReturnsErrorOnNonSerializableObject(t *testing.T) {
	serializer := NewBaseSerializer[NonSerializableObject]()

	object := NonSerializableObject{Function: func() string { return "hello" }}
	_, err := serializer.Serialize(object)

	testtools.AssertErrorIsNotNil(t, err)
}
