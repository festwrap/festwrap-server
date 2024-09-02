package setlistfm

import (
	"encoding/json"
	errors "festwrap/internal/serialization/errors"
	"festwrap/internal/setlist"
)

type SetlistFMDeserializer struct {
	minimumSongs int
}

func (s *SetlistFMDeserializer) SetMinimumSongs(minimumSongs int) {
	s.minimumSongs = minimumSongs
}

func (s *SetlistFMDeserializer) Deserialize(setlist []byte) (*setlist.Setlist, error) {
	var response SetlistFMResponse
	err := json.Unmarshal(setlist, &response)
	if err != nil {
		return nil, errors.NewDeserializationError(err.Error())
	}

	return s.findSetlistWithMinSongs(response), nil
}

func (s *SetlistFMDeserializer) findSetlistWithMinSongs(response SetlistFMResponse) *setlist.Setlist {
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

func NewSetlistFMDeserializer() SetlistFMDeserializer {
	return SetlistFMDeserializer{minimumSongs: 1}
}
