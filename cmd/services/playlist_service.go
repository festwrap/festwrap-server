package playlist

import (
	"context"

	"festwrap/internal/playlist"
)

type CreationStatus int

const (
	Success CreationStatus = iota
	PartialFailure
)

type PlaylistCreation struct {
	PlaylistId string
	Status     CreationStatus
}

type PlaylistService interface {
	CreatePlaylistWithArtists(
		ctx context.Context,
		playlist playlist.PlaylistDetails,
		artists []string,
	) (PlaylistCreation, error)
}
