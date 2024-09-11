package spotify

import (
	"festwrap/internal/artist"
	"festwrap/internal/artist/errors"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"fmt"
	"net/url"
)

type SpotifyArtistRepository struct {
	accessToken  string
	host         string
	deserializer serialization.Deserializer[[]artist.Artist]
	httpSender   httpsender.HTTPRequestSender
}

func (r *SpotifyArtistRepository) SetDeserializer(deserializer serialization.Deserializer[[]artist.Artist]) {
	r.deserializer = deserializer
}

func (r *SpotifyArtistRepository) SearchArtist(name string, limit int) (*[]artist.Artist, error) {

	httpOptions := r.createSetlistHttpOptions(name, limit)
	responseBody, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return nil, errors.NewCannotRetrieveArtistsError(err.Error())
	}

	artists, err := r.deserializer.Deserialize(*responseBody)
	if err != nil {
		return nil, errors.NewCannotRetrieveArtistsError(err.Error())
	}

	return artists, nil
}

func (r *SpotifyArtistRepository) createSetlistHttpOptions(artist string, limit int) httpsender.HTTPRequestOptions {
	httpOptions := httpsender.NewHTTPRequestOptions(r.getSearchUrl(artist, limit), httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", r.accessToken)},
	)
	return httpOptions
}

func (r *SpotifyArtistRepository) getSearchUrl(artistName string, limit int) string {
	queryParams := url.Values{}
	queryParams.Set("type", "artist")
	queryParams.Set("q", fmt.Sprintf("artist:%s", artistName))
	queryParams.Set("limit", fmt.Sprint(limit))
	return fmt.Sprintf("https://%s/v1/search?%s", r.host, queryParams.Encode())
}

func NewSpotifyArtistRepository(accessToken string, httpSender httpsender.HTTPRequestSender) *SpotifyArtistRepository {
	deserializer := NewSpotifyArtistDeserializer()
	return &SpotifyArtistRepository{
		accessToken:  accessToken,
		host:         "api.spotify.com",
		deserializer: &deserializer,
		httpSender:   httpSender,
	}
}
