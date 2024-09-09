package spotify

import (
	"path/filepath"
	"reflect"
	"testing"

	"festwrap/internal/serialization/errors"
	"festwrap/internal/song"
	"festwrap/internal/testtools"
)

func expectedSongs(t *testing.T) *[]song.Song {
	t.Helper()

	songs := []song.Song{
		song.NewSong("spotify:track:4rH1kFLYW0b28UNRyn7dK3"),
		song.NewSong("spotify:track:2pl1Yo26URVBFQRrJXvyuX"),
		song.NewSong("spotify:track:5R0JuZYJxvTKAUnbtoGBXt"),
		song.NewSong("spotify:track:4Qk98QWnZBSJtnqHW7GGpZ"),
		song.NewSong("spotify:track:2I0KOx4fOuS9BV613HLOZN"),
	}
	return &songs
}

func deserializeResponse(t *testing.T, deserializer SpotifySongsDeserializer) []song.Song {
	response := testtools.LoadTestDataOrError(
		t, filepath.Join(testtools.GetParentDir(t), "testdata", "search_song_response.json"),
	)
	result, err := deserializer.Deserialize(response)
	if err != nil {
		t.Fatalf("Found error while parsing: %v", err)
	}
	return *result
}

func TestSongRetrieved(t *testing.T) {
	deserializer := NewSpotifySongsDeserializer()

	result := deserializeResponse(t, deserializer)

	expected := expectedSongs(t)
	testtools.AssertEqual(t, result, *expected)
}

func TestReturnsErrorWhenResponseIsNotJson(t *testing.T) {
	deserializer := NewSpotifySongsDeserializer()

	_, err := deserializer.Deserialize([]byte("{some: non_json}"))

	if _, ok := err.(*errors.DeserializationError); !ok {
		t.Errorf("Expected deserialization error, found %v", reflect.TypeOf(err))
	}

}
