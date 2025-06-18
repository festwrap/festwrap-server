package spotify

import (
	"context"
	"errors"

	types "festwrap/internal"
	"festwrap/internal/artist"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"fmt"
	"net/url"
)

type SpotifyArtistRepository struct {
	tokenKey     types.ContextKey
	host         string
	deserializer serialization.Deserializer[spotifyResponse]
	httpSender   httpsender.HTTPRequestSender
}

func NewSpotifyArtistRepository(httpSender httpsender.HTTPRequestSender) SpotifyArtistRepository {
	deserializer := serialization.NewJsonDeserializer[spotifyResponse]()
	return SpotifyArtistRepository{
		tokenKey:     "token",
		host:         "api.spotify.com",
		deserializer: deserializer,
		httpSender:   httpSender,
	}
}

func (r *SpotifyArtistRepository) SetTokenKey(key types.ContextKey) {
	r.tokenKey = key
}

func (r *SpotifyArtistRepository) SetDeserializer(deserializer serialization.Deserializer[spotifyResponse]) {
	r.deserializer = deserializer
}

func (r *SpotifyArtistRepository) SearchArtist(ctx context.Context, name string, limit int) ([]artist.Artist, error) {
	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return nil, errors.New("could not retrieve token from context")
	}

	httpOptions := r.createSetlistHttpOptions(name, limit, token)
	responseBody, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var response spotifyResponse
	err = r.deserializer.Deserialize(*responseBody, &response)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return response.GetArtists(), nil
}

func (r *SpotifyArtistRepository) createSetlistHttpOptions(
	artist string,
	limit int,
	token string,
) httpsender.HTTPRequestOptions {
	httpOptions := httpsender.NewHTTPRequestOptions(r.getSearchUrl(artist, limit), httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
	)
	return httpOptions
}

func (r *SpotifyArtistRepository) getSearchUrl(artistName string, limit int) string {
	queryParams := url.Values{}
	queryParams.Set("type", "artist")
	queryParams.Set("q", artistName)
	queryParams.Set("limit", fmt.Sprint(limit))
	return fmt.Sprintf("https://%s/v1/search?%s", r.host, queryParams.Encode())
}
