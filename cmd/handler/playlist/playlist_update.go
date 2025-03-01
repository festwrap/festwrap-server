package playlist

type playlistArtist struct {
	Name string `json:"name"`
}

type playlistUpdate struct {
	Artists []playlistArtist `json:"artists"`
}
