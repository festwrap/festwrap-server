package spotify

import (
	"festwrap/internal/song"
	"festwrap/internal/testtools"
	"testing"
)

func TestSpotifySongsSerializer(t *testing.T) {
	songs := []song.Song{song.NewSong("first_uri"), song.NewSong("second_uri")}
	serializer := SpotifySongsSerializer{}

	actual, err := serializer.Serialize(songs)

	expected := []byte("{\"uris\":[\"first_uri\",\"second_uri\"]}")
	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, actual, expected)
}
