package spotify

type spotifySong struct {
	Uri string `json:"uri"`
}

type spotifyTracks struct {
	Songs []spotifySong `json:"items"`
}

type spotifyResponse struct {
	Tracks spotifyTracks `json:"tracks"`
}
