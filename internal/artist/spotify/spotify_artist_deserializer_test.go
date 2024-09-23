package spotify

import (
	"path/filepath"
	"reflect"
	"testing"

	"festwrap/internal/artist"
	"festwrap/internal/serialization/errors"
	"festwrap/internal/testtools"
)

func expectedArtists(t *testing.T) *[]artist.Artist {
	t.Helper()

	artists := []artist.Artist{
		artist.NewArtistWithImageUri(
			"The Beatles",
			"https://i.scdn.co/image/ab6761610000f178e9348cc01ff5d55971b22433",
		),
		artist.NewArtistWithImageUri(
			"The Beatles Tribute Band",
			"https://i.scdn.co/image/ab67616d00004851a53d58fac4e46d5264adc122",
		),
		artist.NewArtistWithImageUri(
			"The Beatles Recovered Band",
			"https://i.scdn.co/image/ab67616d000048518439bc625fdf63ed50bbb594",
		),
		artist.NewArtistWithImageUri(
			"The Beatles Greatest Hits Performed By The Frank Berman Band",
			"https://i.scdn.co/image/ab67616d00004851f903d75acdce7727b3c4aa2c",
		),
		artist.NewArtistWithImageUri(
			"The Beatles Revival Band",
			"https://i.scdn.co/image/ab67616d0000485128cc0b64d05bf392b52a46d3",
		),
	}
	return &artists
}

func expectedArtistsNoImages(t *testing.T) *[]artist.Artist {
	t.Helper()

	artists := []artist.Artist{
		artist.NewArtistWithImageUri(
			"The Beatles",
			"https://i.scdn.co/image/ab6761610000f178e9348cc01ff5d55971b22433",
		),
		artist.NewArtist("The Beatles Tribute Band"),
		artist.NewArtist("The Beatles Recovered Band"),
		artist.NewArtistWithImageUri(
			"The Beatles Greatest Hits Performed By The Frank Berman Band",
			"https://i.scdn.co/image/ab67616d00004851f903d75acdce7727b3c4aa2c",
		),
		artist.NewArtist("The Beatles Revival Band"),
	}
	return &artists
}

func artistsFilePath(t *testing.T) string {
	return filepath.Join(testtools.GetParentDir(t), "testdata", "spotify_artist_search_response.json")
}

func artistsNoImagesFilePath(t *testing.T) string {
	return filepath.Join(testtools.GetParentDir(t), "testdata", "spotify_artist_search_response_no_images.json")
}

func deserializeResponse(t *testing.T, path string, deserializer SpotifyArtistDeserializer) []artist.Artist {
	response := testtools.LoadTestDataOrError(t, path)
	result, err := deserializer.Deserialize(response)
	if err != nil {
		t.Fatalf("Found error while parsing: %v", err)
	}
	return *result
}

func TestArtistsRetrieved(t *testing.T) {
	deserializer := NewSpotifyArtistDeserializer()

	result := deserializeResponse(t, artistsFilePath(t), deserializer)

	expected := expectedArtists(t)
	testtools.AssertEqual(t, result, *expected)
}

func TestReturnsErrorWhenResponseIsNotJson(t *testing.T) {
	deserializer := NewSpotifyArtistDeserializer()

	_, err := deserializer.Deserialize([]byte("{some: non_json}"))

	if _, ok := err.(*errors.DeserializationError); !ok {
		t.Errorf("Expected deserialization error, found %v", reflect.TypeOf(err))
	}
}

func TestArtistsRetrievedWhenTheyHaveNoImages(t *testing.T) {
	deserializer := NewSpotifyArtistDeserializer()

	result := deserializeResponse(t, artistsNoImagesFilePath(t), deserializer)

	expected := expectedArtistsNoImages(t)
	testtools.AssertEqual(t, result, *expected)
}
