package spotify

import (
	"context"
	"errors"
	"fmt"

	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/playlist"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
)

type SpotifyPlaylistRepository struct {
	songsSerializer            serialization.Serializer[SpotifySongs]
	playlistCreateSerializer   serialization.Serializer[SpotifyPlaylist]
	playlistCreateDeserializer serialization.Deserializer[SpotifyCreatePlaylistResponse]
	userIdKey                  types.ContextKey
	tokenKey                   types.ContextKey
	host                       string
	httpSender                 httpsender.HTTPRequestSender
}

func NewSpotifyPlaylistRepository(httpSender httpsender.HTTPRequestSender) SpotifyPlaylistRepository {
	songSerializer := serialization.NewJsonSerializer[SpotifySongs]()
	playlistCreateSerializer := serialization.NewJsonSerializer[SpotifyPlaylist]()
	playlistCreateDeserializer := serialization.NewJsonDeserializer[SpotifyCreatePlaylistResponse]()
	return SpotifyPlaylistRepository{
		tokenKey:                   "token",
		userIdKey:                  "user_id",
		host:                       "api.spotify.com",
		httpSender:                 httpSender,
		songsSerializer:            &songSerializer,
		playlistCreateSerializer:   &playlistCreateSerializer,
		playlistCreateDeserializer: playlistCreateDeserializer,
	}
}

func (r *SpotifyPlaylistRepository) AddSongs(ctx context.Context, playlistId string, songs []song.Song) error {
	if len(songs) == 0 {
		return errors.New("no songs provided")
	}

	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return errors.New("could not retrieve token from context")
	}

	body, err := r.songsSerializer.Serialize(NewSpotifySongs(songs))
	if err != nil {
		errorMsg := fmt.Sprintf("could not serialize songs: %v", err.Error())
		return errors.New(errorMsg)
	}

	httpOptions := r.addSongsHttpOptions(playlistId, body, token)
	_, err = r.httpSender.Send(httpOptions)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r *SpotifyPlaylistRepository) CreatePlaylist(ctx context.Context, playlist playlist.Playlist) (string, error) {
	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return "", errors.New("could not retrieve token from context")
	}

	userId, ok := ctx.Value(r.userIdKey).(string)
	if !ok {
		return "", errors.New("could not retrieve user id from context")
	}

	body, err := r.playlistCreateSerializer.Serialize(
		SpotifyPlaylist{
			Name:        playlist.Name,
			Description: playlist.Description,
			IsPublic:    playlist.IsPublic,
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not serialize playlist: %v", err.Error())
	}

	httpOptions := r.createPlaylistOptions(userId, body, token)
	response, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return "", errors.New(err.Error())
	}

	var parsedResponse SpotifyCreatePlaylistResponse
	err = r.playlistCreateDeserializer.Deserialize(*response, &parsedResponse)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return parsedResponse.Id, nil
}

func (r *SpotifyPlaylistRepository) SetUserIdKey(key types.ContextKey) {
	r.userIdKey = key
}

func (r *SpotifyPlaylistRepository) SetTokenKey(key types.ContextKey) {
	r.tokenKey = key
}

func (r *SpotifyPlaylistRepository) GetHTTPSender() httpsender.HTTPRequestSender {
	return r.httpSender
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

func (r *SpotifyPlaylistRepository) SetPlaylistCreateSerializer(serializer serialization.Serializer[SpotifyPlaylist]) {
	r.playlistCreateSerializer = serializer
}
