package spotify

import (
	"festwrap/internal/playlist"
)

type FakePlaylistSerializer struct {
	response []byte
	err      error
}

func (s *FakePlaylistSerializer) SetResponse(response []byte) {
	s.response = response
}

func (s *FakePlaylistSerializer) SetError(err error) {
	s.err = err
}

func (s *FakePlaylistSerializer) Serialize(playlist playlist.Playlist) ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}
