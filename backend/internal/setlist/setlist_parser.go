package setlist

type SetlistParser interface {
	Parse(setlist []byte) (*Setlist, error)
}
