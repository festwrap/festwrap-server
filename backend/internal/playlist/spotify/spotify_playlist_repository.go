package spotify

import (
	"fmt"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/playlist/errors"
	"festwrap/internal/song"
)

type SpotifyPlaylistRepository struct {
	songsSerializer SongsSerializer
	accessToken     string
	host            string
	httpSender      httpsender.HTTPRequestSender
}

type SpotifyTrackUris struct {
	Uris []string `json:"uris"`
}

func (r *SpotifyPlaylistRepository) AddSongs(playlistId string, songs []song.Song) error {
	if len(songs) == 0 {
		return errors.NewCannotAddSongsToPlaylistError("No songs provided")
	}

	body, err := r.songsSerializer.Serialize(songs)
	if err != nil {
		errorMsg := fmt.Sprintf("Could not serialize request body: %v", err.Error())
		return errors.NewCannotAddSongsToPlaylistError(errorMsg)
	}

	httpOptions := r.createPlaylistHttpOptions(playlistId, body)
	_, err = r.httpSender.Send(httpOptions)
	if err != nil {
		return errors.NewCannotAddSongsToPlaylistError(err.Error())
	}

	return nil
}

func (r *SpotifyPlaylistRepository) SetSongSerializer(serializer SongsSerializer) {
	r.songsSerializer = serializer
}

func (r *SpotifyPlaylistRepository) createPlaylistHttpOptions(playlistId string, body []byte) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://%s/v1/playlists/%s/tracks", r.host, playlistId)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", r.accessToken),
		"Content-Type":  "application/json",
	}
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	httpOptions.SetBody(body)
	httpOptions.SetHeaders(headers)
	return httpOptions
}

func NewSpotifyPlaylistRepository(
	host string, httpSender httpsender.HTTPRequestSender, accessToken string) SpotifyPlaylistRepository {
	return SpotifyPlaylistRepository{
		accessToken:     accessToken,
		host:            host,
		httpSender:      httpSender,
		songsSerializer: &SpotifySongsSerializer{},
	}
}
