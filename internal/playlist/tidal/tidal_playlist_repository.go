package tidal

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

type TidalPlaylistRepository struct {
	songsSerializer            serialization.Serializer[tidalSongs]
	playlistCreateSerializer   serialization.Serializer[tidalPlaylist]
	playlistCreateDeserializer serialization.Deserializer[tidalCreatePlaylistResponse]
	tokenKey                   types.ContextKey
	host                       string
	httpSender                 httpsender.HTTPRequestSender
}

func NewTidalPlaylistRepository(httpSender httpsender.HTTPRequestSender) TidalPlaylistRepository {
	songSerializer := serialization.NewJsonSerializer[tidalSongs]()
	playlistCreateSerializer := serialization.NewJsonSerializer[tidalPlaylist]()
	playlistCreateDeserializer := serialization.NewJsonDeserializer[tidalCreatePlaylistResponse]()
	return TidalPlaylistRepository{
		tokenKey:                   "token",
		host:                       "openapi.tidal.com",
		httpSender:                 httpSender,
		songsSerializer:            &songSerializer,
		playlistCreateSerializer:   &playlistCreateSerializer,
		playlistCreateDeserializer: playlistCreateDeserializer,
	}
}

func (r *TidalPlaylistRepository) AddSongs(ctx context.Context, playlistId string, songs []song.Song) error {
	if len(songs) == 0 {
		return errors.New("no songs provided")
	}

	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return errors.New("could not retrieve token from context while adding songs")
	}

	body, err := r.songsSerializer.Serialize(NewTidalSongs(songs))
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

func (r *TidalPlaylistRepository) CreatePlaylist(ctx context.Context, playlist playlist.PlaylistDetails) (string, error) {
	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return "", errors.New("could not retrieve token from context when creating playlist")
	}

	body, err := r.playlistCreateSerializer.Serialize(
		tidalPlaylist{
			Data: tidalPlaylistData{
				Name:        playlist.Name,
				Description: playlist.Description,
				AccessType:  "PUBLIC",
			},
			Type: "playlist",
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not serialize playlist: %v", err.Error())
	}

	httpOptions := r.createPlaylistOptions(body, token)
	response, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return "", errors.New(err.Error())
	}

	var parsedResponse tidalCreatePlaylistResponse
	err = r.playlistCreateDeserializer.Deserialize(*response, &parsedResponse)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return parsedResponse.Data.Id, nil
}

func (r *TidalPlaylistRepository) SetTokenKey(key types.ContextKey) {
	r.tokenKey = key
}

func (r *TidalPlaylistRepository) GetHTTPSender() httpsender.HTTPRequestSender {
	return r.httpSender
}

func (r *TidalPlaylistRepository) addSongsHttpOptions(
	playlistId string, body []byte, token string,
) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://%s/v2/playlists/%s/relationships/items", r.host, playlistId)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	httpOptions.SetBody(body)
	httpOptions.SetHeaders(r.getTidalBaseHeaders(token))
	return httpOptions
}

func (r *TidalPlaylistRepository) createPlaylistOptions(
	body []byte, token string,
) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://%s/v2/playlists", r.host)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	httpOptions.SetBody(body)
	httpOptions.SetHeaders(r.getTidalBaseHeaders(token))
	return httpOptions
}

func (r *TidalPlaylistRepository) getTidalBaseHeaders(token string) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}
}

func (r *TidalPlaylistRepository) SetPlaylistCreateSerializer(serializer serialization.Serializer[tidalPlaylist]) {
	r.playlistCreateSerializer = serializer
}
