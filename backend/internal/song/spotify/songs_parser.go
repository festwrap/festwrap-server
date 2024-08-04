package spotify

import "festwrap/internal/song"

type SongsParser interface {
	Parse(songs []byte) ([]song.Song, error)
}
