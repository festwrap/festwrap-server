package setlistfm

import (
	"errors"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/setlist"
	"festwrap/internal/testtools"
	"testing"
)

func defaultArtist() string {
	return "Boysetsfire"
}

func defaultSetlist() setlist.Setlist {
	songs := []setlist.Song{
		setlist.NewSong("Closure"),
		setlist.NewSong("Rookie"),
		setlist.NewSong("The Misery Index"),
		setlist.NewSong("One Match"),
		setlist.NewSong("Requiem"),
	}
	return setlist.NewSetlist(defaultArtist(), songs)
}

func senderResponse() []byte {
	return []byte("some response")
}

func defaultSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	response := senderResponse()
	sender.SetResponse(&response)
	return &sender
}

func defaultDeserializer() serialization.FakeDeserializer[setlist.Setlist] {
	result := serialization.FakeDeserializer[setlist.Setlist]{}
	response := defaultSetlist()
	result.SetResponse(&response)
	return result
}

func expectedHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://api.setlist.fm/rest/1.0/search/setlists?artistName=Boysetsfire"
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{
			"x-api-key": "some_api_key",
			"Accept":    "application/json",
		},
	)
	return options
}

func setlistRepository(sender httpsender.HTTPRequestSender) SetlistFMRepository {
	repository := *NewSetlistFMSetlistRepository("some_api_key", sender)
	deserializer := defaultDeserializer()
	repository.SetDeserializer(&deserializer)
	return repository
}

func TestGetSetlistSenderCalledWithProperOptions(t *testing.T) {
	sender := defaultSender()
	repository := setlistRepository(sender)

	repository.GetSetlist(defaultArtist())

	testtools.AssertEqual(t, sender.GetSendArgs(), expectedHttpOptions())
}

func TestGetSetlistReturnsErrorOnSenderError(t *testing.T) {
	sender := httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test error"))
	repository := setlistRepository(&sender)

	_, err := repository.GetSetlist(defaultArtist())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSetlistDeserializerCalledWithSenderResponse(t *testing.T) {
	repository := setlistRepository(defaultSender())
	deserializer := defaultDeserializer()
	repository.SetDeserializer(&deserializer)

	repository.GetSetlist(defaultArtist())

	testtools.AssertEqual(t, deserializer.GetArgs(), senderResponse())
}

func TestGetSetlistReturnsErrorOnDeserializationError(t *testing.T) {
	repository := setlistRepository(defaultSender())
	deserializer := defaultDeserializer()
	deserializer.SetError(errors.New("test deserialization error"))
	repository.SetDeserializer(&deserializer)

	_, err := repository.GetSetlist(defaultArtist())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSetlistReturnsErrorOnEmptySetlist(t *testing.T) {
	repository := setlistRepository(defaultSender())
	deserializer := defaultDeserializer()
	deserializer.SetResponse(nil)
	repository.SetDeserializer(&deserializer)

	_, err := repository.GetSetlist(defaultArtist())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSetlistReturnsResponseSetlist(t *testing.T) {
	repository := setlistRepository(defaultSender())

	actual, _ := repository.GetSetlist(defaultArtist())

	testtools.AssertEqual(t, *actual, defaultSetlist())
}
