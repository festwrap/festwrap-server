package setlistfm

import (
	"encoding/json"
	errors "festwrap/internal/serialization/errors"
)

type SetlistFMDeserializer struct {
	minimumSongs int
}

func NewSetlistFMDeserializer() SetlistFMDeserializer {
	return SetlistFMDeserializer{minimumSongs: 1}
}

func (s *SetlistFMDeserializer) SetMinimumSongs(minimumSongs int) {
	s.minimumSongs = minimumSongs
}

func (s *SetlistFMDeserializer) Deserialize(setlist []byte) (*setlistFMResponse, error) {
	var response setlistFMResponse
	err := json.Unmarshal(setlist, &response)
	if err != nil {
		return &setlistFMResponse{}, errors.NewDeserializationError(err.Error())
	}

	return &response, nil
}
