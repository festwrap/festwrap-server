package playlist_mocks

import (
	"net/http"

	"festwrap/internal/playlist"

	"github.com/stretchr/testify/mock"
)

type PlaylistUpdateBuilderMock struct {
	mock.Mock
}

func (b *PlaylistUpdateBuilderMock) Build(request *http.Request) (playlist.PlaylistUpdate, error) {
	args := b.Called(request)
	return args.Get(0).(playlist.PlaylistUpdate), args.Error(1)
}
