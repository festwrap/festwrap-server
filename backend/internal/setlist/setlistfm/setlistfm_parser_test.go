package setlistfm

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"festwrap/internal/setlist"
	"festwrap/internal/setlist/errors"
	"festwrap/internal/testtools"
)

func loadResponse(t *testing.T) []byte {
	t.Helper()

	dataPath := filepath.Join(testtools.GetParentDir(t), "testdata", "response.json")
	data, err := os.ReadFile(dataPath)

	if err != nil {
		t.Fatalf("Could not load test response: %v", err)
	}
	return data
}

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

func parseResponse(t *testing.T, parser SetlistFMParser) *setlist.Setlist {
	response := loadResponse(t)
	result, err := parser.Parse(response)
	if err != nil {
		t.Fatalf("Found error while parsing: %v", err)
	}
	return result
}

func TestSetlistRetrieved(t *testing.T) {
	parser := NewSetlistFMParser()

	actual := parseResponse(t, parser)

	expected := expectedSetlist(t)
	testtools.AssertEqual(t, actual, expected)
}

func TestNoSetlistRetrievedWhenMinSongsNotReached(t *testing.T) {
	parser := NewSetlistFMParser()
	parser.SetMinimumSongs(25)

	result := parseResponse(t, parser)

	if result != nil {
		t.Errorf("Expected nil, found %v", result)
	}
}

func TestReturnsErrorWhenResponseIsNotJson(t *testing.T) {
	parser := NewSetlistFMParser()

	_, err := parser.Parse([]byte("{some: non_json}"))

	if _, ok := err.(*errors.CannotParseSetlistError); !ok {
		t.Errorf("Expected parse setlist error, found %v", reflect.TypeOf(err))
	}

}
