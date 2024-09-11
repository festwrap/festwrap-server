package setlistfm

import (
	"path/filepath"
	"reflect"
	"testing"

	"festwrap/internal/serialization/errors"
	"festwrap/internal/setlist"
	"festwrap/internal/testtools"
)

func expectedSetlist(t *testing.T) *setlist.Setlist {
	t.Helper()

	songs := []setlist.Song{
		setlist.NewSong("Walk of Life"),
		setlist.NewSong("Anna"),
		setlist.NewSong("Nice Things"),
		setlist.NewSong("America (You're Freaking Me Out)"),
		setlist.NewSong("The Obituaries"),
		setlist.NewSong("After the Party"),
		setlist.NewSong("Irish Goodbyes"),
		setlist.NewSong("Casey"),
		setlist.NewSong("Layla"),
	}
	setlist := setlist.NewSetlist("The Menzingers", songs)
	return &setlist
}

func deserializeResponse(t *testing.T, deserializer SetlistFMDeserializer) *setlist.Setlist {
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

	expected := expectedSetlist(t)
	testtools.AssertEqual(t, actual, expected)
}

func TestNoSetlistRetrievedWhenMinSongsNotReached(t *testing.T) {
	deserializer := NewSetlistFMDeserializer()
	deserializer.SetMinimumSongs(25)

	result := deserializeResponse(t, deserializer)

	if result != nil {
		t.Errorf("Expected nil, found %v", result)
	}
}

func TestReturnsErrorWhenResponseIsNotJson(t *testing.T) {
	deserializer := NewSetlistFMDeserializer()

	_, err := deserializer.Deserialize([]byte("{some: non_json}"))

	if _, ok := err.(*errors.DeserializationError); !ok {
		t.Errorf("Expected deserialization error, found %v", reflect.TypeOf(err))
	}

}
