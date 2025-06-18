package playlist

import "context"

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
	CreatePlaylistWithArtists(ctx context.Context, playlist Playlist, artists []PlaylistArtist) (PlaylistCreation, error)
}
