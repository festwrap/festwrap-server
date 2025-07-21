package playlist

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	playlistmocks "festwrap/internal/playlist/mocks"
	"festwrap/internal/setlist"
	setlistmocks "festwrap/internal/setlist/mocks"
	"festwrap/internal/song"
	songmocks "festwrap/internal/song/mocks"
	"festwrap/internal/testtools"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	playlistId          = "myPlaylist"
	artistName          = "myArtist"
	playlistName        = "My playlist"
	playlistDescription = "Some playlist"
	isPlaylistPublic    = true
)

type SetlistValue struct {
	value setlist.Setlist
	err   error
}

type SongResult struct {
	value song.Song
	err   error
}

type TestArtist struct {
	name    string
	setlist SetlistValue
	songs   []SongResult
}

func (a *TestArtist) SetSetlistError() {
	a.setlist = SetlistValue{err: errors.New("setlist test error")}
}

func (a *TestArtist) SetAllSongsError() {
	for i := range a.songs {
		a.songs[i].err = errors.New("song test error")
	}
}

func (a *TestArtist) SetFirstSongError() {
	a.songs[0].err = errors.New("song test error")
}

func (a *TestArtist) SetEmptySetlist() {
	a.setlist.value = setlist.NewSetlist(a.name, []setlist.Song{}, "https://empty_setlist")
}

func mainTestCase() []TestArtist {
	return []TestArtist{
		{
			name: "Alexisonfire",
			setlist: SetlistValue{
				value: setlist.NewSetlist(
					"Alexisonfire",
					[]setlist.Song{setlist.NewSong("Crisis"), setlist.NewSong("Accidents")},
					"https://alexisonfire",
				),
				err: nil,
			},
			songs: []SongResult{
				{value: song.NewSong("http://some_url1")},
				{value: song.NewSong("http://some_url2")},
			},
		},
		{
			name: "AFI",
			setlist: SetlistValue{
				value: setlist.NewSetlist("AFI", []setlist.Song{setlist.NewSong("Silver and cold")}, "https://afi"),
				err:   nil,
			},
			songs: []SongResult{{value: song.NewSong("http://some_url3")}},
		},
	}
}

func allSetlistsFailTestCase() []TestArtist {
	testArtists := mainTestCase()
	for i := range testArtists {
		testArtists[i].SetSetlistError()
	}
	return testArtists
}

func allSongsFailTestCase() []TestArtist {
	testArtists := mainTestCase()
	for i := range testArtists {
		testArtists[i].SetAllSongsError()
	}
	return testArtists
}

func allSetlistsEmptyTestCase() []TestArtist {
	testArtists := mainTestCase()
	for i := range testArtists {
		testArtists[i].SetEmptySetlist()
	}
	return testArtists
}

func someSetlistsFailTestCase() []TestArtist {
	testArtists := mainTestCase()
	testArtists[0].SetSetlistError()
	return testArtists
}

func someSongsFailedTestCase() []TestArtist {
	testArtists := mainTestCase()
	testArtists[0].SetFirstSongError()
	return testArtists
}

func someSetlistEmptyTestCase() []TestArtist {
	testArtists := mainTestCase()
	testArtists[0].SetEmptySetlist()
	return testArtists
}

func testArtistNames() []string {
	var names []string
	for _, artist := range mainTestCase() {
		names = append(names, artist.name)
	}
	return names
}

func testContext() context.Context {
	return context.Background()
}

func testPlaylist() playlist.PlaylistDetails {
	return playlist.PlaylistDetails{Name: playlistName, Description: playlistDescription, IsPublic: isPlaylistPublic}
}

func newPlaylistRepositoryMock(artists []TestArtist) *playlistmocks.PlaylistRepositoryMock {
	repository := playlistmocks.NewPlaylistRepositoryMock()
	repository.On("CreatePlaylist", testContext(), testPlaylist()).Return(playlistId, nil)
	for _, artist := range artists {
		var songs []song.Song
		for _, song := range artist.songs {
			if song.err == nil {
				songs = append(songs, song.value)
			}
		}
		songsMatcher := mock.MatchedBy(func(items []song.Song) bool {
			return testtools.HaveSameElements(songs, items)
		})
		repository.On(
			"AddSongs",
			testContext(),
			playlistId,
			songsMatcher,
		).Return(nil)
	}
	return &repository
}

func newSetlistRepositoryMock(artists []TestArtist) *setlistmocks.SetlistRepositoryMock {
	repository := setlistmocks.NewSetlistRepositoryMock()
	for _, artist := range artists {
		repository.On("GetSetlist", artist.name, mock.Anything).Return(artist.setlist.value, artist.setlist.err)
	}
	return &repository
}

func newSongRepositoryMock(artists []TestArtist) *songmocks.SongRepositoryMock {
	repository := songmocks.NewSongRepositoryMock()
	for _, artist := range artists {
		for i, setlistSong := range artist.setlist.value.GetSongs() {
			repository.On(
				"GetSong",
				testContext(),
				artist.name,
				setlistSong.GetTitle(),
			).Return(artist.songs[i].value, artist.songs[i].err)
		}
	}
	return &repository
}

func testSetup(testCase []TestArtist) (
	*playlistmocks.PlaylistRepositoryMock,
	*setlistmocks.SetlistRepositoryMock,
	*songmocks.SongRepositoryMock,
) {
	playlistRepository := newPlaylistRepositoryMock(testCase)
	setlistRepository := newSetlistRepositoryMock(testCase)
	songRepository := newSongRepositoryMock(testCase)
	return playlistRepository, setlistRepository, songRepository
}

func TestCreatePlaylistPerformsExpectedSideEffects(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup(mainTestCase())
	service := NewBasePlaylistService(playlistRepository, setlistRepository, songRepository, logging.NoopLogger{})

	_, err := service.CreatePlaylistWithArtists(testContext(), testPlaylist(), testArtistNames())

	assert.Nil(t, err)
	playlistRepository.AssertExpectations(t)
	setlistRepository.AssertExpectations(t)
}

func TestCreatePlaylistReturnsErrorOnCreateError(t *testing.T) {
	_, setlistRepository, songRepository := testSetup(mainTestCase())
	playlistRepository := playlistmocks.NewPlaylistRepositoryMock()
	playlistRepository.On("CreatePlaylist", testContext(), testPlaylist()).Return("", errors.New("create test error"))
	service := NewBasePlaylistService(
		&playlistRepository, setlistRepository, songRepository, logging.NoopLogger{})

	_, err := service.CreatePlaylistWithArtists(testContext(), testPlaylist(), testArtistNames())

	assert.NotNil(t, err)
}

func TestCreatePlaylistResult(t *testing.T) {
	tests := map[string]struct {
		testCase       []TestArtist
		expectedStatus PlaylistCreation
		expectedError  error
	}{
		"all setlists fail": {
			testCase:       allSetlistsFailTestCase(),
			expectedStatus: PlaylistCreation{},
			expectedError:  fmt.Errorf("all artists failed to be added to playlist %s", playlistId),
		},
		"all songs failed to be added": {
			testCase:       allSongsFailTestCase(),
			expectedStatus: PlaylistCreation{},
			expectedError:  fmt.Errorf("all artists failed to be added to playlist %s", playlistId),
		},
		"all setlists empty": {
			testCase:       allSetlistsEmptyTestCase(),
			expectedStatus: PlaylistCreation{},
			expectedError:  fmt.Errorf("all artists failed to be added to playlist %s", playlistId),
		},
		"some setlists failed": {
			testCase:       someSetlistsFailTestCase(),
			expectedStatus: PlaylistCreation{PlaylistId: playlistId, Status: PartialFailure},
			expectedError:  nil,
		},
		"some setlists empty": {
			testCase:       someSetlistEmptyTestCase(),
			expectedStatus: PlaylistCreation{PlaylistId: playlistId, Status: PartialFailure},
			expectedError:  nil,
		},
		"some songs failed": {
			testCase:       someSongsFailedTestCase(),
			expectedStatus: PlaylistCreation{PlaylistId: playlistId, Status: Success},
			expectedError:  nil,
		},
		"success": {
			testCase:       mainTestCase(),
			expectedStatus: PlaylistCreation{PlaylistId: playlistId, Status: Success},
			expectedError:  nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			playlistRepository, setlistRepository, songRepository := testSetup(test.testCase)
			service := NewBasePlaylistService(
				playlistRepository, setlistRepository, songRepository, logging.NoopLogger{})

			status, err := service.CreatePlaylistWithArtists(testContext(), testPlaylist(), testArtistNames())

			assert.Equal(t, test.expectedStatus, status)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
