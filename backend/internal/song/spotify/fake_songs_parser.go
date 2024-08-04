package spotify

import (
	"festwrap/internal/song"
)

type FakeSongsParser struct {
	response  []song.Song
	err       error
	parseArgs []byte
}

func (s *FakeSongsParser) SetResponse(response []song.Song) {
	s.response = response
}

func (s *FakeSongsParser) SetError(err error) {
	s.err = err
}

func (s *FakeSongsParser) GetParseArgs() []byte {
	return s.parseArgs
}

func (s *FakeSongsParser) Parse(setlist []byte) ([]song.Song, error) {
	s.parseArgs = setlist
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}
