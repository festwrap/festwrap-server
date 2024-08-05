package spotify

import (
	"errors"
	"testing"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/song"
	"festwrap/internal/testtools"
)

func fakeSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	emptyResponse := []byte("")
	sender.SetResponse(&emptyResponse)
	return &sender
}

func errorSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test send error"))
	return &sender
}

func defaultPlaylistId() string {
	return "test_id"
}

func defaultSongs() []song.Song {
	return []song.Song{song.NewSong("uri1"), song.NewSong("uri2")}
}

func defaultBody() []byte {
	return []byte("some body")
}

func defaultSerializer() *FakeSongsSerializer {
	serializer := FakeSongsSerializer{}
	serializer.SetResponse(defaultBody())
	return &serializer
}

func errorSerializer() SongsSerializer {
	serializer := FakeSongsSerializer{}
	serializer.SetError(errors.New("test error"))
	return &serializer
}

func expectedHttpOptions() httpsender.HTTPRequestOptions {
	options := httpsender.NewHTTPRequestOptions("https://spotify.com/v1/playlists/test_id/tracks", httpsender.POST, 201)
	options.SetHeaders(
		map[string]string{
			"Authorization": "Bearer abcdefg12345",
			"Content-Type":  "application/json",
		},
	)
	options.SetBody(defaultBody())
	return options
}

func spotifyPlaylistRepository() SpotifyPlaylistRepository {
	repository := NewSpotifyPlaylistRepository("spotify.com", fakeSender(), "abcdefg12345")
	repository.SetSongSerializer(defaultSerializer())
	return repository
}

func TestAddSongsReturnsErrorWhenNoSongsProvided(t *testing.T) {
	repository := spotifyPlaylistRepository()

	err := repository.AddSongs(defaultPlaylistId(), []song.Song{})

	testtools.AssertErrorIsNotNil(t, err)
}

func TestAddSongsReturnsErrorOnNonSerializableInput(t *testing.T) {
	repository := spotifyPlaylistRepository()
	repository.SetSongSerializer(errorSerializer())

	err := repository.AddSongs(defaultPlaylistId(), defaultSongs())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestAddSongsSendsRequestUsingProperOptions(t *testing.T) {
	sender := fakeSender()
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(sender)

	repository.AddSongs(defaultPlaylistId(), defaultSongs())

	actual := sender.GetSendArgs()
	testtools.AssertEqual(t, actual, expectedHttpOptions())
}

func TestAddSongsReturnsErrorOnSendError(t *testing.T) {
	sender := errorSender()
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(sender)

	err := repository.AddSongs(defaultPlaylistId(), defaultSongs())

	testtools.AssertErrorIsNotNil(t, err)
}
