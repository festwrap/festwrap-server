package event

type PlaylistCreationStatus string
type PlaylistType string

const (
	PLAYLIST_CREATED_OK             PlaylistCreationStatus = "ok"
	PLAYLIST_CREATED_PARTIAL_ERRORS PlaylistCreationStatus = "partial_error"
)

const (
	PLAYLIST_TYPE_SPOTIFY PlaylistType = "spotify"
)

type CreatedPlaylist struct {
	Id      string       `json:"id"`
	Name    string       `json:"name"`
	Artists []string     `json:"artists"`
	Type    PlaylistType `json:"type"`
}

type PlaylistCreatedEvent struct {
	Playlist       CreatedPlaylist        `json:"playlist"`
	CreationStatus PlaylistCreationStatus `json:"status"`
}

func (e PlaylistCreatedEvent) Type() EventType {
	return PlaylistCreated
}
