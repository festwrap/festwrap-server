package playlist

import "context"

type CreationStatus int

const (
	Success CreationStatus = iota
	PartialFailure
	Error
)

type PlaylistCreation struct {
	PlaylistId string
	Status     CreationStatus
}

type PlaylistService interface {
	CreatePlaylistWithArtists(
		ctx context.Context, playlist PlaylistDetails, artists []PlaylistArtist) (PlaylistCreation, error)
}
