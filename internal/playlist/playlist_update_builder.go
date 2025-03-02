package playlist

import (
	"net/http"
)

type PlaylistArtist struct {
	Name string
}

type PlaylistUpdate struct {
	PlaylistId string
	Artists    []PlaylistArtist
}

type PlaylistUpdateBuilder interface {
	Build(request *http.Request) (PlaylistUpdate, error)
}
