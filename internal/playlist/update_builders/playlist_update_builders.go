package update_builders

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"festwrap/internal/playlist"
	"festwrap/internal/serialization"
)

type PlaylistUpdateBuilder interface {
	Build(request *http.Request) (playlist.PlaylistUpdate, error)
}

type ExistingPlaylistUpdateBuilder struct {
	pathId       string
	deserializer serialization.Deserializer[ExistingPlaylistUpdate]
}

type NewPlaylistUpdateBuilder struct {
	playlistService playlist.PlaylistService
	deserializer    serialization.Deserializer[NewPlaylistUpdate]
}

func NewExistingPlaylistUpdateBuilder(pathId string) ExistingPlaylistUpdateBuilder {
	return ExistingPlaylistUpdateBuilder{
		pathId:       pathId,
		deserializer: serialization.NewJsonDeserializer[ExistingPlaylistUpdate](),
	}
}

func (b *ExistingPlaylistUpdateBuilder) Build(request *http.Request) (playlist.PlaylistUpdate, error) {
	playlistId := request.PathValue(b.pathId)
	if playlistId == "" {
		return playlist.PlaylistUpdate{}, fmt.Errorf("could not find playlist id in %s", b.pathId)
	}

	defer request.Body.Close()
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return playlist.PlaylistUpdate{}, errors.New("could read body from request")
	}

	var artists ExistingPlaylistUpdate
	err = b.deserializer.Deserialize(requestBody, &artists)
	if err != nil {
		return playlist.PlaylistUpdate{}, errors.New("failed to deserialize playlist artists: " + err.Error())
	}

	updateArtists := make([]playlist.PlaylistArtist, len(artists.Artists))
	for i, artist := range artists.Artists {
		updateArtists[i] = playlist.PlaylistArtist{Name: artist.Name}
	}
	update := playlist.PlaylistUpdate{PlaylistId: playlistId, Artists: updateArtists}
	return update, nil
}

func NewNewPlaylistUpdateBuilder(playlistService playlist.PlaylistService) NewPlaylistUpdateBuilder {
	return NewPlaylistUpdateBuilder{
		playlistService: playlistService,
		deserializer:    serialization.NewJsonDeserializer[NewPlaylistUpdate](),
	}
}

func (b *NewPlaylistUpdateBuilder) Build(request *http.Request) (playlist.PlaylistUpdate, error) {
	defer request.Body.Close()
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return playlist.PlaylistUpdate{}, errors.New("could read body from request")
	}

	var update NewPlaylistUpdate
	err = b.deserializer.Deserialize(requestBody, &update)
	if err != nil {
		return playlist.PlaylistUpdate{}, errors.New("failed to deserialize playlist information: " + err.Error())
	}

	playlistId, err := b.playlistService.CreatePlaylist(
		request.Context(),
		update.Playlist.toPlaylist(),
	)
	if err != nil {
		return playlist.PlaylistUpdate{}, errors.New("could not create playlist")
	}

	playlistArtists := make([]playlist.PlaylistArtist, len(update.Artists))
	for i, artist := range update.Artists {
		playlistArtists[i] = playlist.PlaylistArtist{Name: artist.Name}
	}
	return playlist.PlaylistUpdate{PlaylistId: playlistId, Artists: playlistArtists}, nil
}
