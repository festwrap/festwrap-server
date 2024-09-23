package setlist

type SetlistRepository interface {
	GetSetlist(artist string) (*Setlist, error)
}
