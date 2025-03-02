package update_builders

type PlaylistArtist struct {
	Name string `json:"name"`
}

type PlaylistArtists struct {
	Artists []PlaylistArtist `json:"artists"`
}
