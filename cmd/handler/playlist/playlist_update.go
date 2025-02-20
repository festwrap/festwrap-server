package playlist

type updateArtists struct {
	Name string `json:"name"`
}

type playlistUpdate struct {
	Artists []updateArtists `json:"artists"`
}
