package playlist

import (
	"festwrap/internal/setlist"
)

type PlaylistMetadata struct {
	Name     string
	Id       string
	Href     string
	ImageURI string
}

type Playlist struct {
	PlaylistMetadata
	setlists []setlist.Setlist
}
