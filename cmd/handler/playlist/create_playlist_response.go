package playlist

type CreatedPlaylist struct {
	Id string `json:"id"`
}

type CreatePlaylistResponse struct {
	Playlist CreatedPlaylist `json:"playlist"`
}
