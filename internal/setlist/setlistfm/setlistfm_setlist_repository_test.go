package setlistfm

import (
	"errors"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/setlist"
	"festwrap/internal/testtools"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	artist = "The Menzingers"
)

func defaultMinSongs() int {
	return 3
}

func responseBody(t *testing.T) []byte {
	path := filepath.Join(testtools.GetParentDir(t), "testdata", "response.json")
	return testtools.LoadTestDataOrError(t, path)
}

func emptyResponseBody(t *testing.T) []byte {
	path := filepath.Join(testtools.GetParentDir(t), "testdata", "no_setlists_response.json")
	return testtools.LoadTestDataOrError(t, path)
}

func defaultSender(t *testing.T) *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	response := responseBody(t)
	sender.SetResponse(&response)
	return &sender
}

func expectedHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://api.setlist.fm/rest/1.0/search/setlists?artistName=The+Menzingers"
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{
			"x-api-key": "some_api_key",
			"Accept":    "application/json",
		},
	)
	return options
}

func expectedSetlist() *setlist.Setlist {
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

func TestGetSetlistSenderCalledWithProperOptions(t *testing.T) {
	sender := defaultSender(t)
	repository := NewSetlistFMSetlistRepository("some_api_key", sender)

	repository.GetSetlist(artist, defaultMinSongs())

	assert.Equal(t, sender.GetSendArgs(), expectedHttpOptions())
}

func TestGetSetlistReturnsErrorOnSenderError(t *testing.T) {
	sender := httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test error"))
	repository := NewSetlistFMSetlistRepository("some_api_key", &sender)

	_, err := repository.GetSetlist(artist, defaultMinSongs())

	assert.NotNil(t, err)
}

func TestGetSetlistReturnsErrorOnDeserializationError(t *testing.T) {
	sender := httpsender.FakeHTTPSender{}
	invalidResponse := []byte("{bad response}")
	sender.SetResponse(&invalidResponse)
	repository := NewSetlistFMSetlistRepository("some_api_key", &sender)

	_, err := repository.GetSetlist(artist, defaultMinSongs())

	assert.NotNil(t, err)
}

func TestGetSetlistReturnsErrorIfNoSetlistFound(t *testing.T) {
	sender := httpsender.FakeHTTPSender{}
	response := emptyResponseBody(t)
	sender.SetResponse(&response)
	repository := NewSetlistFMSetlistRepository("some_api_key", &sender)

	_, err := repository.GetSetlist(artist, defaultMinSongs())

	assert.NotNil(t, err)
}

func TestGetSetlistReturnsSetlist(t *testing.T) {
	repository := NewSetlistFMSetlistRepository("some_api_key", defaultSender(t))

	actual, _ := repository.GetSetlist(artist, defaultMinSongs())

	assert.Equal(t, actual, expectedSetlist())
}

func TestGetSetlistRetrievesErrorWhenMinSongsNotReached(t *testing.T) {
	repository := NewSetlistFMSetlistRepository("some_api_key", defaultSender(t))

	_, err := repository.GetSetlist(artist, 50)

	assert.NotNil(t, err)
}
