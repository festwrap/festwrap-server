package spotify

import (
	"festwrap/internal/song"
)

type SongsSerializer interface {
	Serialize(songs []song.Song) ([]byte, error)
}
