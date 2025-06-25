package setlist_mocks

import (
	"festwrap/internal/setlist"

	"github.com/stretchr/testify/mock"
)

type SetlistRepositoryMock struct {
	mock.Mock
}

func NewSetlistRepositoryMock() SetlistRepositoryMock {
	return SetlistRepositoryMock{}
}

func (s *SetlistRepositoryMock) GetSetlist(artist string, minSongs int) (setlist.Setlist, error) {
	args := s.Called(artist, minSongs)
	return args.Get(0).(setlist.Setlist), args.Error(1)
}
