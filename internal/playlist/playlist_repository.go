package playlist

import "festwrap/internal/song"

type PlaylistRepository interface {
	CreatePlaylist(userId string, playlist Playlist) error
	AddSongs(playlistId string, songs []song.Song) error
}
