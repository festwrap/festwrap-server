package spotify

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"festwrap/internal/song"
	"festwrap/internal/song/errors"
	"festwrap/internal/testtools"
)

func loadResponse(t *testing.T) []byte {
	t.Helper()

	dataPath := filepath.Join(testtools.GetParentDir(t), "testdata", "search_song_response.json")
	data, err := os.ReadFile(dataPath)

	if err != nil {
		t.Fatalf("Could not load test response: %v", err)
	}
	return data
}

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

func parseResponse(t *testing.T, parser SpotifySongsParser) []song.Song {
	response := loadResponse(t)
	result, err := parser.Parse(response)
	if err != nil {
		t.Fatalf("Found error while parsing: %v", err)
	}
	return result
}

func TestSongRetrieved(t *testing.T) {
	parser := NewSpotifySongsParser()

	result := parseResponse(t, parser)

	expected := expectedSongs(t)
	testtools.AssertEqual(t, result, *expected)
}

func TestReturnsErrorWhenResponseIsNotJson(t *testing.T) {
	parser := NewSpotifySongsParser()

	_, err := parser.Parse([]byte("{some: non_json}"))

	if _, ok := err.(*errors.CannotParseSongsError); !ok {
		t.Errorf("Expected parse song error, found %v", reflect.TypeOf(err))
	}

}
