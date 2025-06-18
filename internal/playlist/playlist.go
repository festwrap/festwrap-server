package playlist

type PlaylistDetails struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

type Playlist struct {
	PlaylistDetails
	Id string `json:"id"`
}
