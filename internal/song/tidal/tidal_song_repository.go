package tidal

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
)

type TidalSongRepository struct {
	tokenKey     types.ContextKey
	host         string
	countryCode  string
	httpSender   httpsender.HTTPRequestSender
	deserializer serialization.Deserializer[tidalResponse]
}

func NewTidalSongRepository(httpSender httpsender.HTTPRequestSender) *TidalSongRepository {
	return &TidalSongRepository{
		tokenKey:     "token",
		host:         "openapi.tidal.com",
		countryCode:  "ES",
		httpSender:   httpSender,
		deserializer: serialization.NewJsonDeserializer[tidalResponse](),
	}
}

func (r *TidalSongRepository) GetSong(ctx context.Context, artist string, title string) (song.Song, error) {
	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return song.Song{}, errors.New("could not retrieve token from context when retrieving song")
	}

	httpOptions := r.createSongHttpOptions(artist, title, token)
	responseBody, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return song.Song{}, errors.New(err.Error())
	}

	var response tidalResponse
	err = r.deserializer.Deserialize(*responseBody, &response)
	if err != nil {
		return song.Song{}, errors.New(err.Error())
	}

	if len(response.Results) == 0 {
		return song.Song{}, fmt.Errorf("no songs found for song %s (%s)", title, artist)
	}

	// We assume the first result is the most trusted one
	result := song.NewSong(response.Results[0].Id)
	return result, nil
}

func (r *TidalSongRepository) SetDeserializer(deserializer serialization.Deserializer[tidalResponse]) {
	r.deserializer = deserializer
}

func (r *TidalSongRepository) createSongHttpOptions(
	artist string,
	title string,
	token string,
) httpsender.HTTPRequestOptions {
	httpOptions := httpsender.NewHTTPRequestOptions(r.getSongFullUrl(artist, title), httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
	)
	return httpOptions
}

func (r *TidalSongRepository) getSongFullUrl(artist string, title string) string {
	searchQuery := fmt.Sprintf("%s %s", artist, title)
	encodedQuery := url.PathEscape(searchQuery)
	queryParams := url.Values{}
	queryParams.Set("countryCode", r.countryCode)
	queryParams.Set("include", "tracks")
	return fmt.Sprintf("https://%s/v2/searchResults?%s&%s", r.host, encodedQuery, queryParams.Encode())
}

func (r *TidalSongRepository) SetTokenKey(key types.ContextKey) {
	r.tokenKey = key
}
