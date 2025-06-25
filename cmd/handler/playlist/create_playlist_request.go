package playlist

type PlaylistArtist struct {
	Name string `json:"name"`
}

type NewPlaylist struct {
	Name string `json:"name"`
}

type NewPlaylistRequest struct {
	Playlist NewPlaylist      `json:"playlist"`
	Artists  []PlaylistArtist `json:"artists"`
}

func (r NewPlaylistRequest) GetArtistNames() []string {
	var names []string
	for _, artist := range r.Artists {
		names = append(names, artist.Name)
	}
	return names
}
