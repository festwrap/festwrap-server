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
	deserializer serialization.Deserializer[PlaylistArtists]
}

func NewExistingPlaylistUpdateBuilder(pathId string) ExistingPlaylistUpdateBuilder {
	return ExistingPlaylistUpdateBuilder{
		pathId:       pathId,
		deserializer: serialization.NewJsonDeserializer[PlaylistArtists](),
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

	var artists PlaylistArtists
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

func (b *ExistingPlaylistUpdateBuilder) SetDeserializer(deserializer serialization.Deserializer[PlaylistArtists]) {
	b.deserializer = deserializer
}
