package spotify

import (
	"errors"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
	"festwrap/internal/testtools"
	"testing"
)

func defaultArtist() string {
	return "Movements"
}

func defaultTitle() string {
	return "Daylily"
}

func expectedHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://api.spotify.com/v1/search?q=artist%3AMovements+track%3ADaylily&type=track"
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{"Authorization": "Bearer some_token"},
	)
	return options
}

func defaultSenderResponse() *[]byte {
	response := []byte("some body")
	return &response
}

func defaultSongs() []song.Song {
	return []song.Song{song.NewSong("some uri"), song.NewSong("another uri")}
}

func defaultSender() *httpsender.FakeHTTPSender {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetResponse(defaultSenderResponse())
	return sender
}

func defaultDeserializer() *serialization.FakeDeserializer[[]song.Song] {
	deserializer := &serialization.FakeDeserializer[[]song.Song]{}
	response := defaultSongs()
	deserializer.SetResponse(&response)
	return deserializer
}

func spotifySongRepository(sender httpsender.HTTPRequestSender) SpotifySongRepository {
	repository := NewSpotifySongRepository("some_token", sender)
	repository.SetDeserializer(defaultDeserializer())
	return *repository
}

func TestGetSongSendsRequestWithProperOptions(t *testing.T) {
	sender := defaultSender()
	repository := spotifySongRepository(sender)

	_, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, sender.GetSendArgs(), expectedHttpOptions())
}

func TestGetSongReturnsErrorOnSendError(t *testing.T) {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test error"))
	repository := spotifySongRepository(sender)

	_, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSongCallsDeserializeWithSendResponseBody(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	deserializer := defaultDeserializer()
	repository.SetDeserializer(deserializer)

	_, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, deserializer.GetArgs(), *defaultSenderResponse())
}

func TestGetSongReturnsErrorOnResponseBodyDeserializationError(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	deserializer := defaultDeserializer()
	deserializer.SetError(errors.New("test error"))
	repository.SetDeserializer(deserializer)

	_, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSongReturnsErrorIfNoSongsFound(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	deserializer := defaultDeserializer()
	deserializer.SetResponse(&[]song.Song{})
	repository.SetDeserializer(deserializer)

	_, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSongReturnsFirstSongFound(t *testing.T) {
	repository := spotifySongRepository(defaultSender())

	song, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, *song, defaultSongs()[0])
}
