package spotify

import (
	"context"
	"fmt"
	"net/url"

	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
	"festwrap/internal/song/errors"
)

type SpotifySongRepository struct {
	tokenKey     types.ContextKey
	host         string
	httpSender   httpsender.HTTPRequestSender
	deserializer serialization.Deserializer[spotifyResponse]
}

func NewSpotifySongRepository(httpSender httpsender.HTTPRequestSender) *SpotifySongRepository {
	return &SpotifySongRepository{
		tokenKey:     "token",
		host:         "api.spotify.com",
		httpSender:   httpSender,
		deserializer: serialization.NewJsonDeserializer[spotifyResponse](),
	}
}

func (r *SpotifySongRepository) GetSong(ctx context.Context, artist string, title string) (*song.Song, error) {
	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return nil, errors.NewCannotRetrieveSongError("Could not retrieve token from context")
	}

	httpOptions := r.createSongHttpOptions(artist, title, token)
	responseBody, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return nil, errors.NewCannotRetrieveSongError(err.Error())
	}

	var response spotifyResponse
	err = r.deserializer.Deserialize(*responseBody, &response)
	if err != nil {
		return nil, errors.NewCannotRetrieveSongError(err.Error())
	}

	if len(response.Tracks.Songs) == 0 {
		errorMsg := fmt.Sprintf("No songs found for song %s (%s)", title, artist)
		return nil, errors.NewCannotRetrieveSongError(errorMsg)
	}

	// We assume the first result is the most trusted one
	result := song.NewSong(response.Tracks.Songs[0].Uri)
	return &result, nil
}

func (r *SpotifySongRepository) SetDeserializer(deserializer serialization.Deserializer[spotifyResponse]) {
	r.deserializer = deserializer
}

func (r *SpotifySongRepository) createSongHttpOptions(
	artist string,
	title string,
	token string,
) httpsender.HTTPRequestOptions {
	httpOptions := httpsender.NewHTTPRequestOptions(r.getSetlistFullUrl(artist, title), httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
	)
	return httpOptions
}

func (r *SpotifySongRepository) getSetlistFullUrl(artist string, title string) string {
	queryParams := url.Values{}
	queryParams.Set("q", fmt.Sprintf("artist:%s track:%s", artist, title))
	queryParams.Set("type", "track")
	setlistPath := "v1/search"
	return fmt.Sprintf("https://%s/%s?%s", r.host, setlistPath, queryParams.Encode())
}

func (r *SpotifySongRepository) SetTokenKey(key types.ContextKey) {
	r.tokenKey = key
}
