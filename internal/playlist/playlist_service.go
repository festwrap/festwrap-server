package playlist

import "context"

type PlaylistService interface {
	CreatePlaylist(ctx context.Context, playlist PlaylistDetails) (string, error)
	AddSetlist(ctx context.Context, playlistId string, artist string) error
}
