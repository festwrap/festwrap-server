package playlist

type PlaylistService interface {
	AddSetlist(playlistId string, artist string)
}
