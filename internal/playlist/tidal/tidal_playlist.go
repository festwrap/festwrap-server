package tidal

type tidalPlaylistData struct {
	AccessType  string `json:"accessType"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type tidalPlaylist struct {
	Data tidalPlaylistData `json:"data"`
	Type string            `json:"type"`
}
