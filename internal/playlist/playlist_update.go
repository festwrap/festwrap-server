package playlist

type PlaylistArtist struct {
	Name string `json:"name"`
}

type PlaylistArtists struct {
	Artists []PlaylistArtist `json:"artists"`
}

type PlaylistUpdate struct {
	PlaylistId string
	Artists    []PlaylistArtist
}
