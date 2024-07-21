package spotify

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"festwrap/internal/song"
	"festwrap/internal/song/errors"
)

type SpotifySongRepositoryConfig struct {
	Host        string
	AccessToken string
}

type SpotifySongRepository struct {
	client *http.Client
	parser *SpotifySongsParser
	config SpotifySongRepositoryConfig
}

func (s *SpotifySongRepository) GetSong(artist string, title string) (*song.Song, error) {

	request, err := s.createSongSearchRequest(artist, title)
	if err != nil {
		errorMsg := fmt.Sprintf("Cannot song search request for %s: %s", s.config.Host, err.Error())
		return nil, errors.NewCannotRetrieveSongError(errorMsg)
	}

	response, err := s.client.Do(request)
	if err != nil {
		errorMsg := fmt.Sprintf(
			"Cannot retrieve song search response for %s: %s", s.config.Host, err.Error(),
		)
		return nil, errors.NewCannotRetrieveSongError(errorMsg)
	}

	if response.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf(
			"Spotify returned %d when trying to retrieve song %s (%s)", response.StatusCode, title, artist,
		)
		return nil, errors.NewCannotRetrieveSongError(errorMsg)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		errorMsg := fmt.Sprintf("Error reading song search response %s: ", err.Error())
		return nil, errors.NewCannotRetrieveSongError(errorMsg)
	}

	songs, err := s.parser.Parse(body)
	if err != nil {
		return nil, err
	}

	allSongs := *songs
	if len(allSongs) == 0 {
		errorMsg := fmt.Sprintf("No songs found for song %s (%s)", title, artist)
		return nil, errors.NewCannotRetrieveSongError(errorMsg)
	}

	// We assume the first result is the most trusted one
	return &allSongs[0], nil
}

func (s *SpotifySongRepository) createSongSearchRequest(artist string, title string) (*http.Request, error) {
	request, err := http.NewRequest("GET", s.getSetlistFullUrl(artist, title), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.config.AccessToken))
	return request, nil
}

func (s *SpotifySongRepository) getSetlistFullUrl(artist string, title string) string {
	queryParams := url.Values{}
	queryParams.Set("q", fmt.Sprintf("artist:%s track:%s", artist, title))
	queryParams.Set("type", "track")
	setlistPath := "v1/search"
	return fmt.Sprintf("https://%s/%s?%s", s.config.Host, setlistPath, queryParams.Encode())
}

func NewSpotifySongRepository(
	client *http.Client,
	config SpotifySongRepositoryConfig,
	parser *SpotifySongsParser,
) *SpotifySongRepository {
	return &SpotifySongRepository{client: client, config: config, parser: parser}
}
