package setlistfm

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	httpsender "festwrap/internal/http/sender"
	httpsendermocks "festwrap/internal/http/sender/mocks"
	"festwrap/internal/setlist"
	"festwrap/internal/testtools"

	"github.com/stretchr/testify/assert"
)

const (
	setlistFMApiKey = "someApiKey"
	artist          = "The Menzingers"
	minSongs        = 3
)

func responseBody(t *testing.T) *[]byte {
	path := filepath.Join(testtools.GetParentDir(t), "testdata", "response.json")
	response := testtools.LoadTestDataOrError(t, path)
	return &response
}

func emptyResponseBody(t *testing.T) *[]byte {
	path := filepath.Join(testtools.GetParentDir(t), "testdata", "no_setlists_response.json")
	response := testtools.LoadTestDataOrError(t, path)
	return &response
}

func sender(t *testing.T) httpsender.HTTPRequestSender {
	sender := httpsendermocks.HTTPSenderMock{}
	sender.On("Send", getSetlistHttpOptions(1)).Return(responseBody(t), nil)
	return &sender
}

func getSetlistHttpOptions(page int) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://api.setlist.fm/rest/1.0/search/setlists?artistName=The+Menzingers&p=%d", page)
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{
			"x-api-key": setlistFMApiKey,
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
	sender := sender(t).(*httpsendermocks.HTTPSenderMock)
	repository := NewSetlistFMSetlistRepository(setlistFMApiKey, sender)

	repository.GetSetlist(artist, minSongs)

	sender.AssertExpectations(t)
}

func TestGetSetlistReturnsErrorOnSenderError(t *testing.T) {
	sender := httpsendermocks.HTTPSenderMock{}
	sender.On("Send", getSetlistHttpOptions(1)).Return(nil, errors.New("test error"))
	repository := NewSetlistFMSetlistRepository(setlistFMApiKey, &sender)

	_, err := repository.GetSetlist(artist, minSongs)

	assert.NotNil(t, err)
}

func TestGetSetlistReturnsErrorOnDeserializationError(t *testing.T) {
	sender := httpsendermocks.HTTPSenderMock{}
	invalidResponse := []byte("{bad response}")
	sender.On("Send", getSetlistHttpOptions(1)).Return(&invalidResponse, nil)
	repository := NewSetlistFMSetlistRepository(setlistFMApiKey, &sender)

	_, err := repository.GetSetlist(artist, minSongs)

	assert.NotNil(t, err)
}

func TestGetSetlistReturnsErrorIfNoSetlistFound(t *testing.T) {
	sender := httpsendermocks.HTTPSenderMock{}
	sender.On("Send", getSetlistHttpOptions(1)).Return(emptyResponseBody(t), nil)
	repository := NewSetlistFMSetlistRepository(setlistFMApiKey, &sender)

	_, err := repository.GetSetlist(artist, minSongs)

	assert.NotNil(t, err)
}

func TestGetSetlistReturnsSetlist(t *testing.T) {
	repository := NewSetlistFMSetlistRepository(setlistFMApiKey, sender(t))

	actual, _ := repository.GetSetlist(artist, minSongs)

	assert.Equal(t, expectedSetlist(), actual)
}

func TestGetSetlistRetrievesErrorWhenMinSongsNotReached(t *testing.T) {
	repository := NewSetlistFMSetlistRepository(setlistFMApiKey, sender(t))

	_, err := repository.GetSetlist(artist, 50)

	assert.NotNil(t, err)
}

func TestGetSetlistReturnsResultsFromNextPageIfFirstHasNoResults(t *testing.T) {
	multiPageSender := httpsendermocks.HTTPSenderMock{}
	multiPageSender.On("Send", getSetlistHttpOptions(1)).Return(emptyResponseBody(t), nil)
	multiPageSender.On("Send", getSetlistHttpOptions(2)).Return(responseBody(t), nil)
	repository := NewSetlistFMSetlistRepository(setlistFMApiKey, &multiPageSender)
	repository.SetMaxPages(3)

	actual, err := repository.GetSetlist(artist, minSongs)

	assert.Equal(t, expectedSetlist(), actual)
	assert.Nil(t, err)
}
