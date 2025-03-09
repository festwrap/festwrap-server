package setlistfm

import (
	"fmt"
	"net/url"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/setlist"
	"festwrap/internal/setlist/errors"
)

type SetlistFMRepository struct {
	host         string
	apiKey       string
	deserializer serialization.Deserializer[setlistFMResponse]
	httpSender   httpsender.HTTPRequestSender
	maxPages     int
}

func NewSetlistFMSetlistRepository(apiKey string, httpSender httpsender.HTTPRequestSender) *SetlistFMRepository {
	deserializer := serialization.NewJsonDeserializer[setlistFMResponse]()
	return &SetlistFMRepository{
		host:         "api.setlist.fm",
		apiKey:       apiKey,
		deserializer: &deserializer,
		httpSender:   httpSender,
		maxPages:     1,
	}
}

func (r *SetlistFMRepository) SetDeserializer(deserializer serialization.Deserializer[setlistFMResponse]) {
	r.deserializer = deserializer
}

func (r *SetlistFMRepository) GetSetlist(artist string, minSongs int) (*setlist.Setlist, error) {

	page := 1
	var setlist *setlist.Setlist
	var err error

	for page <= r.maxPages {
		setlist, err = r.getFirstSetlistFromPage(artist, page, minSongs)
		resultOrErrorFound := setlist != nil || err != nil
		if resultOrErrorFound {
			break
		} else {
			page += 1
		}
	}

	if setlist == nil {
		errorMsg := fmt.Sprintf("Could not find setlist for artist %s", artist)
		return nil, errors.NewCannotRetrieveSetlistError(errorMsg)
	} else {
		return setlist, nil
	}
}

func (r *SetlistFMRepository) getFirstSetlistFromPage(artist string, page int, minSongs int) (*setlist.Setlist, error) {
	httpOptions := r.createSetlistHttpOptions(artist, page)
	responseBody, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return nil, errors.NewCannotRetrieveSetlistError(err.Error())
	}

	var response setlistFMResponse
	err = r.deserializer.Deserialize(*responseBody, &response)
	if err != nil {
		return nil, errors.NewCannotRetrieveSetlistError(err.Error())
	}

	setlist := response.findSetlistWithMinSongs(minSongs)
	return setlist, nil
}

func (r *SetlistFMRepository) createSetlistHttpOptions(artist string, page int) httpsender.HTTPRequestOptions {
	url := r.getSetlistFullUrl(artist, page)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{
			"x-api-key": r.apiKey,
			"Accept":    "application/json",
		},
	)
	return httpOptions
}

func (r *SetlistFMRepository) getSetlistFullUrl(artist string, page int) string {
	queryParams := url.Values{}
	queryParams.Set("artistName", artist)
	queryParams.Set("p", fmt.Sprint(page))
	setlistPath := "rest/1.0/search/setlists"
	return fmt.Sprintf("https://%s/%s?%s", r.host, setlistPath, queryParams.Encode())
}

func (r *SetlistFMRepository) SetMaxPages(maxPages int) {
	r.maxPages = maxPages
}
