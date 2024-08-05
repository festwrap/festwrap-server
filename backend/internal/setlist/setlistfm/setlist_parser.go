package setlistfm

import "festwrap/internal/setlist"

type SetlistParser interface {
	Parse(setlist []byte) (*setlist.Setlist, error)
}
