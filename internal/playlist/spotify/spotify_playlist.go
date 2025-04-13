package spotify

type SpotifyPlaylist struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"public"`
}
