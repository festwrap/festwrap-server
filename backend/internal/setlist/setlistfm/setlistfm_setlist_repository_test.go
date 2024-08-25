package setlistfm

import (
	"errors"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/setlist"
	"festwrap/internal/testtools"
	"testing"
)

func defaultArtist() string {
	return "Boysetsfire"
}

func defaultParsedSetlist() setlist.Setlist {
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

func defaultParser() FakeSetlistParser {
	parser := FakeSetlistParser{}
	response := defaultParsedSetlist()
	parser.SetReponse(&response)
	return parser
}

func expectedHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://setlistm.com/rest/1.0/search/setlists?artistName=Boysetsfire"
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
	repository := *NewSetlistFMSetlistRepository("setlistm.com", "some_api_key", sender)
	parser := defaultParser()
	repository.SetParser(&parser)
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

func TestGetSetlistParserCalledWithSenderResponse(t *testing.T) {
	repository := setlistRepository(defaultSender())
	parser := defaultParser()
	repository.SetParser(&parser)

	repository.GetSetlist(defaultArtist())

	testtools.AssertEqual(t, parser.GetParseArgs(), senderResponse())
}

func TestGetSetlistReturnsErrorOnParserError(t *testing.T) {
	repository := setlistRepository(defaultSender())
	parser := FakeSetlistParser{}
	parser.SetError(errors.New("test parser error"))
	repository.SetParser(&parser)

	_, err := repository.GetSetlist(defaultArtist())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSetlistReturnsErrorOnEmptyParse(t *testing.T) {
	repository := setlistRepository(defaultSender())
	parser := FakeSetlistParser{}
	parser.SetReponse(nil)
	repository.SetParser(&parser)

	_, err := repository.GetSetlist(defaultArtist())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetSetlistReturnsParsedSetlist(t *testing.T) {
	repository := setlistRepository(defaultSender())

	actual, _ := repository.GetSetlist(defaultArtist())

	testtools.AssertEqual(t, *actual, defaultParsedSetlist())
}
