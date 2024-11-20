package spotify

import (
	"context"
	"fmt"

	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/playlist"
	"festwrap/internal/playlist/errors"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
)

type SpotifyPlaylistRepository struct {
	songsSerializer    serialization.Serializer[SpotifySongs]
	playlistSerializer serialization.Serializer[SpotifyPlaylist]
	tokenKey           types.ContextKey
	host               string
	httpSender         httpsender.HTTPRequestSender
}

func NewSpotifyPlaylistRepository(httpSender httpsender.HTTPRequestSender) SpotifyPlaylistRepository {
	playlistSerializer := serialization.NewJsonSerializer[SpotifyPlaylist]()
	songSerializer := serialization.NewJsonSerializer[SpotifySongs]()
	return SpotifyPlaylistRepository{
		tokenKey:           "token",
		host:               "api.spotify.com",
		httpSender:         httpSender,
		songsSerializer:    &songSerializer,
		playlistSerializer: &playlistSerializer,
	}
}

func (r *SpotifyPlaylistRepository) AddSongs(ctx context.Context, playlistId string, songs []song.Song) error {
	if len(songs) == 0 {
		return errors.NewCannotAddSongsToPlaylistError("no songs provided")
	}

	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return errors.NewCannotAddSongsToPlaylistError("Could not retrieve token from context")
	}

	body, err := r.songsSerializer.Serialize(NewSpotifySongs(songs))
	if err != nil {
		errorMsg := fmt.Sprintf("could not serialize songs: %v", err.Error())
		return errors.NewCannotAddSongsToPlaylistError(errorMsg)
	}

	httpOptions := r.addSongsHttpOptions(playlistId, body, token)
	_, err = r.httpSender.Send(httpOptions)
	if err != nil {
		return errors.NewCannotAddSongsToPlaylistError(err.Error())
	}

	return nil
}

func (r *SpotifyPlaylistRepository) CreatePlaylist(ctx context.Context, userId string, playlist playlist.Playlist) error {
	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return errors.NewCannotCreatePlaylistError("Could not retrieve token from context")
	}

	body, err := r.playlistSerializer.Serialize(
		SpotifyPlaylist{
			Name:        playlist.Name,
			Description: playlist.Description,
			IsPublic:    playlist.IsPublic,
		},
	)
	if err != nil {
		errorMsg := fmt.Sprintf("could not serialize playlist: %v", err.Error())
		return errors.NewCannotCreatePlaylistError(errorMsg)
	}

	httpOptions := r.createPlaylistOptions(userId, body, token)
	_, err = r.httpSender.Send(httpOptions)
	if err != nil {
		return errors.NewCannotCreatePlaylistError(err.Error())
	}

	return nil
}

func (r *SpotifyPlaylistRepository) SetTokenKey(key types.ContextKey) {
	r.tokenKey = key
}

func (r *SpotifyPlaylistRepository) SetHTTPSender(httpSender httpsender.HTTPRequestSender) {
	r.httpSender = httpSender
}

func (r *SpotifyPlaylistRepository) SetSongSerializer(serializer serialization.Serializer[SpotifySongs]) {
	r.songsSerializer = serializer
}

func (r *SpotifyPlaylistRepository) SetPlaylistSerializer(serializer serialization.Serializer[SpotifyPlaylist]) {
	r.playlistSerializer = serializer
}

func (r *SpotifyPlaylistRepository) addSongsHttpOptions(
	playlistId string, body []byte, token string,
) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://%s/v1/playlists/%s/tracks", r.host, playlistId)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	httpOptions.SetBody(body)
	httpOptions.SetHeaders(r.GetSpotifyBaseHeaders(token))
	return httpOptions
}

func (r *SpotifyPlaylistRepository) createPlaylistOptions(
	userId string, body []byte, token string,
) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://%s/v1/users/%s/playlists", r.host, userId)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	httpOptions.SetBody(body)
	httpOptions.SetHeaders(r.GetSpotifyBaseHeaders(token))
	return httpOptions
}

func (r *SpotifyPlaylistRepository) GetSpotifyBaseHeaders(token string) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}
}
