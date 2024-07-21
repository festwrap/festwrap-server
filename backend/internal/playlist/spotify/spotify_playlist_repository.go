package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"festwrap/internal/playlist/errors"
	"festwrap/internal/song"
)

type SpotifyPlaylistRepository struct {
	host        string
	accessToken string
	client      *http.Client
}

type SpotifyTrackUris struct {
	Uris []string `json:"uris"`
}

func (r *SpotifyPlaylistRepository) AddSongs(playlistId string, songs []song.Song) error {
	body, err := createRequestBody(songs)
	if err != nil {
		errorMsg := fmt.Sprintf("Could not serialize request body: %v", err.Error())
		return errors.NewCannotAddSongsToPlaylistError(errorMsg)
	}

	url := fmt.Sprintf("https://%s/v1/playlists/%s/tracks", r.host, playlistId)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.accessToken))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		errorMsg := fmt.Sprintf("Could not create request: %v", err.Error())
		return errors.NewCannotAddSongsToPlaylistError(errorMsg)
	}

	response, err := r.client.Do(request)
	if err != nil {
		errorMsg := fmt.Sprintf(
			"Error sending request to add %v to playlist %s: %v", songs, playlistId, err.Error(),
		)
		return errors.NewCannotAddSongsToPlaylistError(errorMsg)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		errorMsg := fmt.Sprintf(
			"Could not add songs %v to playlist %s, responsed with %d", songs, playlistId, response.StatusCode,
		)
		return errors.NewCannotAddSongsToPlaylistError(errorMsg)
	}

	return nil
}

func createRequestBody(songs []song.Song) ([]byte, error) {
	songUris := []string{}
	for _, currentSong := range songs {
		songUris = append(songUris, currentSong.GetUri())
	}
	return json.Marshal(SpotifyTrackUris{Uris: songUris})
}

func NewSpotifyPlaylistRepository(
	host string, client *http.Client, accessToken string) SpotifyPlaylistRepository {
	return SpotifyPlaylistRepository{host: host, client: client, accessToken: accessToken}
}
