package playlist

import (
	"context"
	"festwrap/internal/song"
)

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, playlist PlaylistDetails) (string, error)
	AddSongs(ctx context.Context, playlistId string, songs []song.Song) error
}
