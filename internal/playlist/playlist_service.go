package playlist

import "context"

type PlaylistService interface {
	CreatePlaylist(ctx context.Context, playlist Playlist) error
	AddSetlist(ctx context.Context, playlistId string, artist string) error
}
