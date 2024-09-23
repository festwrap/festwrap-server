package spotify

import (
	"fmt"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/playlist"
	"festwrap/internal/playlist/errors"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
)

type SpotifyPlaylistRepository struct {
	songsSerializer    serialization.Serializer[SongList]
	playlistSerializer serialization.Serializer[playlist.Playlist]
	accessToken        string
	host               string
	httpSender         httpsender.HTTPRequestSender
}

func NewSpotifyPlaylistRepository(
	httpSender httpsender.HTTPRequestSender, accessToken string) SpotifyPlaylistRepository {
	return SpotifyPlaylistRepository{
		accessToken:        accessToken,
		host:               "api.spotify.com",
		httpSender:         httpSender,
		songsSerializer:    &SpotifySongsSerializer{},
		playlistSerializer: &SpotifyPlaylistSerializer{},
	}
}

func (r *SpotifyPlaylistRepository) AddSongs(playlistId string, songs []song.Song) error {
	if len(songs) == 0 {
		return errors.NewCannotAddSongsToPlaylistError("no songs provided")
	}

	body, err := r.songsSerializer.Serialize(SongList{songs: songs})
	if err != nil {
		errorMsg := fmt.Sprintf("could not serialize songs: %v", err.Error())
		return errors.NewCannotAddSongsToPlaylistError(errorMsg)
	}

	httpOptions := r.addSongsHttpOptions(playlistId, body)
	_, err = r.httpSender.Send(httpOptions)
	if err != nil {
		return errors.NewCannotAddSongsToPlaylistError(err.Error())
	}

	return nil
}

func (r *SpotifyPlaylistRepository) CreatePlaylist(userId string, playlist playlist.Playlist) error {
	body, err := r.playlistSerializer.Serialize(playlist)
	if err != nil {
		errorMsg := fmt.Sprintf("could not serialize playlist: %v", err.Error())
		return errors.NewCannotCreatePlaylistError(errorMsg)
	}

	httpOptions := r.createPlaylistOptions(userId, body)
	_, err = r.httpSender.Send(httpOptions)
	if err != nil {
		return errors.NewCannotCreatePlaylistError(err.Error())
	}

	return nil
}

func (r *SpotifyPlaylistRepository) SetHTTPSender(httpSender httpsender.HTTPRequestSender) {
	r.httpSender = httpSender
}

func (r *SpotifyPlaylistRepository) SetSongSerializer(serializer serialization.Serializer[SongList]) {
	r.songsSerializer = serializer
}

func (r *SpotifyPlaylistRepository) SetPlaylistSerializer(serializer serialization.Serializer[playlist.Playlist]) {
	r.playlistSerializer = serializer
}

func (r *SpotifyPlaylistRepository) addSongsHttpOptions(playlistId string, body []byte) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://%s/v1/playlists/%s/tracks", r.host, playlistId)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	httpOptions.SetBody(body)
	httpOptions.SetHeaders(r.GetSpotifyBaseHeaders())
	return httpOptions
}

func (r *SpotifyPlaylistRepository) createPlaylistOptions(userId string, body []byte) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://%s/v1/users/%s/playlists", r.host, userId)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	httpOptions.SetBody(body)
	httpOptions.SetHeaders(r.GetSpotifyBaseHeaders())
	return httpOptions
}

func (r *SpotifyPlaylistRepository) GetSpotifyBaseHeaders() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", r.accessToken),
		"Content-Type":  "application/json",
	}
}
