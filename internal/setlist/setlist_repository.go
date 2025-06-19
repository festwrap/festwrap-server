package setlist

type SetlistRepository interface {
	GetSetlist(artist string, minSongs int) (Setlist, error)
}
