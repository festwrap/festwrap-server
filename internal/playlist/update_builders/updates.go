package update_builders

import "festwrap/internal/playlist"

type PlaylistArtist struct {
	Name string `json:"name"`
}

type ExistingPlaylistUpdate struct {
	Artists []PlaylistArtist `json:"artists"`
}

type NewPlaylist struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

func (p NewPlaylist) toPlaylist() playlist.PlaylistDetails {
	return playlist.PlaylistDetails{
		Name:        p.Name,
		Description: p.Description,
		IsPublic:    p.IsPublic,
	}
}

type NewPlaylistUpdate struct {
	ExistingPlaylistUpdate
	Playlist NewPlaylist `json:"playlist"`
}
