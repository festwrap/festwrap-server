package song_mocks

import (
	"context"

	"festwrap/internal/song"

	"github.com/stretchr/testify/mock"
)

type SongRepositoryMock struct {
	mock.Mock
}

func NewSongRepositoryMock() SongRepositoryMock {
	return SongRepositoryMock{}
}

func (s *SongRepositoryMock) GetSong(ctx context.Context, artist string, title string) (song.Song, error) {
	args := s.Called(ctx, artist, title)
	return args.Get(0).(song.Song), args.Error(1)
}
