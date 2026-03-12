package tidal

type tidalCreatePlaylistResponseData struct {
	Id string `json:"id"`
}

type tidalCreatePlaylistResponse struct {
	Data tidalCreatePlaylistResponseData `json:"data"`
}
