package serialization

type Object struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type NonSerializableObject struct {
	Function func() string
}

func serializableObjectBytes() []byte {
	return []byte(`{"name":"myname","value":10}`)
}

func serializableObject() Object {
	return Object{Name: "myname", Value: 10}
}

func nonSerializableObject() NonSerializableObject {
	return NonSerializableObject{Function: func() string { return "hello" }}
}
