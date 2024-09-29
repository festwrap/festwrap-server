package setlistfm

import (
	"path/filepath"
	"reflect"
	"testing"

	"festwrap/internal/serialization/errors"
	"festwrap/internal/testtools"
)

func expectedResponse() *setlistFMResponse {
	set := setlistfmSet{
		[]setlistfmSong{
			{"Walk of Life"},
			{"Anna"},
			{"Nice Things"},
			{"America (You're Freaking Me Out)"},
			{"The Obituaries"},
			{"After the Party"},
		},
	}
	encore := setlistfmSet{
		[]setlistfmSong{
			{"Irish Goodbyes"},
			{"Casey"},
			{"Layla"},
		},
	}
	artist := setlistfmArtist{Name: "The Menzingers"}
	return &setlistFMResponse{
		Body: []setlistFMSetlist{
			{
				Artist: artist,
				Sets:   setlistFMSets{Sets: []setlistfmSet{}},
			},
			{
				Artist: artist,
				Sets:   setlistFMSets{Sets: []setlistfmSet{}},
			},
			{
				Artist: artist,
				Sets:   setlistFMSets{Sets: []setlistfmSet{set, encore}},
			},
		},
	}
}

func deserializeResponse(t *testing.T, deserializer SetlistFMDeserializer) *setlistFMResponse {
	response := testtools.LoadTestDataOrError(t, filepath.Join(testtools.GetParentDir(t), "testdata", "response.json"))
	result, err := deserializer.Deserialize(response)
	if err != nil {
		t.Fatalf("Found error while parsing: %v", err)
	}
	return result
}

func TestSetlistRetrieved(t *testing.T) {
	deserializer := NewSetlistFMDeserializer()

	actual := deserializeResponse(t, deserializer)

	expected := expectedResponse()
	testtools.AssertEqual(t, *actual, *expected)
}

func TestReturnsErrorWhenResponseIsNotJson(t *testing.T) {
	deserializer := NewSetlistFMDeserializer()

	_, err := deserializer.Deserialize([]byte("{some: non_json}"))

	if _, ok := err.(*errors.DeserializationError); !ok {
		t.Errorf("Expected deserialization error, found %v", reflect.TypeOf(err))
	}

}
