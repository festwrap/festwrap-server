package playlist_mocks

import (
	"context"
	"festwrap/internal/playlist"
	"festwrap/internal/song"

	"github.com/stretchr/testify/mock"
)

type PlaylistRepositoryMock struct {
	mock.Mock
}

func NewPlaylistRepositoryMock() PlaylistRepositoryMock {
	return PlaylistRepositoryMock{}
}

func (s *PlaylistRepositoryMock) CreatePlaylist(ctx context.Context, playlistInput playlist.PlaylistDetails) (string, error) {
	args := s.Called(ctx, playlistInput)
	return args.String(0), args.Error(1)
}

func (s *PlaylistRepositoryMock) AddSongs(ctx context.Context, playlistId string, songs []song.Song) error {
	return s.Called(ctx, playlistId, songs).Error(0)
}
