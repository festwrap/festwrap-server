package spotify

import (
	"festwrap/internal/song"
)

type FakeSongsSerializer struct {
	response []byte
	err      error
}

func (s *FakeSongsSerializer) SetResponse(response []byte) {
	s.response = response
}

func (s *FakeSongsSerializer) SetError(err error) {
	s.err = err
}

func (s *FakeSongsSerializer) Serialize(songs []song.Song) ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}
