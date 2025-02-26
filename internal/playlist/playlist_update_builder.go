package playlist

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"festwrap/internal/serialization"
)

type PlaylistUpdateBuilder interface {
	Build(request *http.Request) (PlaylistUpdate, error)
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

func (b *ExistingPlaylistUpdateBuilder) Build(request *http.Request) (PlaylistUpdate, error) {
	playlistId := request.PathValue(b.pathId)
	if playlistId == "" {
		return PlaylistUpdate{}, fmt.Errorf("could not find playlist id in %s", b.pathId)
	}

	defer request.Body.Close()
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return PlaylistUpdate{}, errors.New("could read body from request")
	}

	var artists PlaylistArtists
	err = b.deserializer.Deserialize(requestBody, &artists)
	if err != nil {
		return PlaylistUpdate{}, errors.New("failed to deserialize playlist artists: " + err.Error())
	}

	update := PlaylistUpdate{PlaylistId: playlistId, Artists: artists.Artists}
	return update, nil
}

func (b *ExistingPlaylistUpdateBuilder) SetDeserializer(deserializer serialization.Deserializer[PlaylistArtists]) {
	b.deserializer = deserializer
}
