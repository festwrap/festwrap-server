package spotify

import (
	"encoding/json"

	"festwrap/internal/song"
	"festwrap/internal/song/errors"
)

type SpotifySongsParser struct{}

func (s *SpotifySongsParser) Parse(setlist []byte) (*[]song.Song, error) {
	var parsedResponse SpotifyResponse
	err := json.Unmarshal(setlist, &parsedResponse)
	if err != nil {
		return nil, errors.NewCannotParseSongsError(err.Error())
	}

	result := []song.Song{}
	for _, currentSong := range parsedResponse.Tracks.Songs {
		result = append(result, song.NewSong(currentSong.Uri))
	}
	return &result, nil
}

func NewSpotifySongsParser() SpotifySongsParser {
	return SpotifySongsParser{}
}
