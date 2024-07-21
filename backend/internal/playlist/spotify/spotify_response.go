package spotify

type SpotifyImage struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

type SpotifyPlaylist struct {
	Name   string         `json:"name"`
	Id     string         `json:"id"`
	Href   string         `json:"href"`
	Images []SpotifyImage `json:"images"`
}

func (p *SpotifyPlaylist) GetFirstImage() *SpotifyImage {
	if len(p.Images) == 0 {
		return nil
	}
	return &p.Images[0]
}

type SpotifyResponse struct {
	Total     int               `json:"total"`
	Playlists []SpotifyPlaylist `json:"items"`
}
