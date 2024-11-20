package spotify

import (
	"fmt"
	"net/url"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
	"festwrap/internal/song/errors"
)

type SpotifySongRepository struct {
	accessToken  string
	host         string
	httpSender   httpsender.HTTPRequestSender
	deserializer serialization.Deserializer[spotifyResponse]
}

func NewSpotifySongRepository(
	accessToken string,
	httpSender httpsender.HTTPRequestSender,
) *SpotifySongRepository {
	return &SpotifySongRepository{
		accessToken:  accessToken,
		host:         "api.spotify.com",
		httpSender:   httpSender,
		deserializer: serialization.NewJsonDeserializer[spotifyResponse](),
	}
}

func (r *SpotifySongRepository) GetSong(artist string, title string) (*song.Song, error) {
	httpOptions := r.createSongHttpOptions(artist, title)
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

func (r *SpotifySongRepository) createSongHttpOptions(artist string, title string) httpsender.HTTPRequestOptions {
	httpOptions := httpsender.NewHTTPRequestOptions(r.getSetlistFullUrl(artist, title), httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", r.accessToken)},
	)
	return httpOptions
}

func (r *SpotifySongRepository) getSetlistFullUrl(artist string, title string) string {
	queryParams := url.Values{}
	queryParams.Set("q", fmt.Sprintf("+artist:%s+track:%s", artist, title))
	queryParams.Set("type", "track")
	setlistPath := "v1/search"
	return fmt.Sprintf("https://%s/%s?%s", r.host, setlistPath, queryParams.Encode())
}
