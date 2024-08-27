package spotify

import "festwrap/internal/playlist"

type PlaylistSerializer interface {
	Serialize(playlist playlist.Playlist) ([]byte, error)
}
