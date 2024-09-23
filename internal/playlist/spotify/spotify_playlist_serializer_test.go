package spotify

import (
	"festwrap/internal/playlist"
	"festwrap/internal/testtools"
	"testing"
)

func TestSpotifyPlaylistSerializer(t *testing.T) {
	playlist := playlist.Playlist{Name: "my-playlist", Description: "my-description", IsPublic: true}
	serializer := SpotifyPlaylistSerializer{}

	actual, err := serializer.Serialize(playlist)

	expected := []byte("{\"name\":\"my-playlist\",\"description\":\"my-description\",\"is_public\":true}")
	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, actual, expected)
}
