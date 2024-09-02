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
	deserializer serialization.Deserializer[setlist.Setlist]
	httpSender   httpsender.HTTPRequestSender
}

func (r *SetlistFMRepository) SetDeserializer(deserializer serialization.Deserializer[setlist.Setlist]) {
	r.deserializer = deserializer
}

func (r *SetlistFMRepository) GetSetlist(artist string) (*setlist.Setlist, error) {

	httpOptions := r.createSetlistHttpOptions(artist)
	responseBody, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return nil, errors.NewCannotRetrieveSetlistError(err.Error())
	}

	setlist, err := r.deserializer.Deserialize(*responseBody)
	if err != nil {
		return nil, errors.NewCannotRetrieveSetlistError(err.Error())
	}
	if setlist == nil {
		// TODO: if no valid setlist found, we should check for the next page
		// TODO: probable a good idea to move the min songs filter to repository
		// TODO: By doing so, we keep deserializer logic simpler
		errorMsg := fmt.Sprintf("Could not find setlist for artist %s", artist)
		return nil, errors.NewCannotRetrieveSetlistError(errorMsg)
	}

	return setlist, nil
}

func (r *SetlistFMRepository) createSetlistHttpOptions(artist string) httpsender.HTTPRequestOptions {
	httpOptions := httpsender.NewHTTPRequestOptions(r.getSetlistFullUrl(artist), httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{
			"x-api-key": r.apiKey,
			"Accept":    "application/json",
		},
	)
	return httpOptions
}

func (r *SetlistFMRepository) getSetlistFullUrl(artist string) string {
	queryParams := url.Values{}
	queryParams.Set("artistName", artist)
	setlistPath := "rest/1.0/search/setlists"
	return fmt.Sprintf("https://%s/%s?%s", r.host, setlistPath, queryParams.Encode())
}

func NewSetlistFMSetlistRepository(apiKey string, httpSender httpsender.HTTPRequestSender) *SetlistFMRepository {
	deserializer := NewSetlistFMDeserializer()
	return &SetlistFMRepository{
		host:         "api.setlist.fm",
		apiKey:       apiKey,
		deserializer: &deserializer,
		httpSender:   httpSender,
	}
}
