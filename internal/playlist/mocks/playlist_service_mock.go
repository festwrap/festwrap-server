package playlist_mocks

import (
	"context"
	"festwrap/internal/playlist"

	"github.com/stretchr/testify/mock"
)

type PlaylistServiceMock struct {
	mock.Mock
}

func NewPlaylistServiceMock() PlaylistServiceMock {
	return PlaylistServiceMock{}
}

func (s *PlaylistServiceMock) CreatePlaylist(ctx context.Context, playlistInput playlist.Playlist) (string, error) {
	args := s.Called(ctx, playlistInput)
	return args.String(0), args.Error(0)
}

func (s *PlaylistServiceMock) AddSetlist(ctx context.Context, playlistId string, artist string) error {
	return s.Called(ctx, playlistId, artist).Error(0)
}
