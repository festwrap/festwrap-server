package playlist_mocks

import (
	"context"
	"festwrap/internal/playlist"

	services "festwrap/cmd/services"

	"github.com/stretchr/testify/mock"
)

type PlaylistServiceMock struct {
	mock.Mock
}

func NewPlaylistServiceMock() PlaylistServiceMock {
	return PlaylistServiceMock{}
}

func (s *PlaylistServiceMock) CreatePlaylistWithArtists(
	ctx context.Context,
	playlist playlist.PlaylistDetails,
	artists []string,
) (services.PlaylistCreation, error) {
	args := s.Called(ctx, playlist, artists)
	return args.Get(0).(services.PlaylistCreation), args.Error(1)
}
