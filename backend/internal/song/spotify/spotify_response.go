package spotify

type SpotifySong struct {
	Uri string `json:"uri"`
}

type SpotifyTracks struct {
	Songs []SpotifySong `json:"items"`
}

type SpotifyResponse struct {
	Tracks SpotifyTracks `json:"tracks"`
}
