package setlistfm

import (
	"encoding/json"
	"festwrap/internal/setlist"
	"festwrap/internal/setlist/errors"
)

type SetlistFMParser struct {
	minimumSongs int
}

func (s *SetlistFMParser) SetMinimumSongs(minimumSongs int) {
	s.minimumSongs = minimumSongs
}

func (s *SetlistFMParser) Parse(setlist []byte) (*setlist.Setlist, error) {
	var parsedResponse SetlistFMResponse
	err := json.Unmarshal(setlist, &parsedResponse)
	if err != nil {
		return nil, errors.NewCannotParseSetlistError(err.Error())
	}

	return s.findSetlistWithMinSongs(parsedResponse), nil
}

func (s *SetlistFMParser) findSetlistWithMinSongs(response SetlistFMResponse) *setlist.Setlist {
	var result *setlist.Setlist
	for _, set := range response.Body {
		currentSetlist := setlist.NewSetlist(set.Artist.Name, set.GetSongs())
		if len(currentSetlist.GetSongs()) >= s.minimumSongs {
			result = &currentSetlist
			break
		}
	}

	return result
}

func NewSetlistFMParser() SetlistFMParser {
	return SetlistFMParser{minimumSongs: 1}
}
