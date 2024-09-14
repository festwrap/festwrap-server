package spotify

import (
	"encoding/json"

	"festwrap/internal/serialization/errors"
	"festwrap/internal/song"
)

type SpotifySongsDeserializer struct{}

func NewSpotifySongsDeserializer() SpotifySongsDeserializer {
	return SpotifySongsDeserializer{}
}

func (s *SpotifySongsDeserializer) Deserialize(setlist []byte) (*[]song.Song, error) {
	var response SpotifyResponse
	err := json.Unmarshal(setlist, &response)
	if err != nil {
		return nil, errors.NewDeserializationError(err.Error())
	}

	result := []song.Song{}
	for _, currentSong := range response.Tracks.Songs {
		result = append(result, song.NewSong(currentSong.Uri))
	}
	return &result, nil
}
