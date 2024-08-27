package spotify

import (
	"encoding/json"
	"festwrap/internal/playlist"
)

type SpotifyPlaylist struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type SpotifyPlaylistSerializer struct{}

func (s *SpotifyPlaylistSerializer) Serialize(playlist playlist.Playlist) ([]byte, error) {
	return json.Marshal(SpotifyPlaylist{Name: playlist.Name, Description: playlist.Description, IsPublic: playlist.IsPublic})
}
