package spotify

import (
	"errors"
	httpsender "festwrap/internal/http/sender"
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
	url := "https://spotify.com/v1/search?q=artist%3AMovements+track%3ADaylily&type=track"
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

func defaultParsedSongs() []song.Song {
	return []song.Song{song.NewSong("some uri"), song.NewSong("another uri")}
}

func defaultSender() *httpsender.FakeHTTPSender {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetResponse(defaultSenderResponse())
	return sender
}

func defaultParser() *FakeSongsParser {
	parser := &FakeSongsParser{}
	parser.SetResponse(defaultParsedSongs())
	return parser
}

func spotifySongRepository(sender httpsender.HTTPRequestSender) SpotifySongRepository {
	repository := NewSpotifySongRepository("some_token", "spotify.com", sender)
	repository.SetParser(defaultParser())
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

func TestGetSongCallsParserWithSendResponseBody(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	parser := defaultParser()
	repository.SetParser(parser)

	_, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, parser.GetParseArgs(), *defaultSenderResponse())
}

func TestGetSongReturnsErrorOnResponseBodyParseError(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	parser := &FakeSongsParser{}
	parser.SetError(errors.New("test error"))
	repository.SetParser(parser)

	_, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSongReturnsErrorIfNoSongsParsed(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	parser := &FakeSongsParser{}
	parser.SetResponse([]song.Song{})
	repository.SetParser(parser)

	_, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSongReturnsFirstSongParsed(t *testing.T) {
	repository := spotifySongRepository(defaultSender())

	song, err := repository.GetSong(defaultArtist(), defaultTitle())

	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, *song, defaultParsedSongs()[0])
}
