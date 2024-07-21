package playlist

import (
	"festwrap/internal/song"
)

type PlaylistRepository interface {
	AddSongs(playlistId string, songs []song.Song) error
}
