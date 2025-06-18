package playlist

import (
	"context"
	"errors"
	"testing"

	"festwrap/internal/setlist"
	"festwrap/internal/song"
	"festwrap/internal/testtools"

	"github.com/stretchr/testify/assert"
)

const (
	playlistId          = "myPlaylist"
	artistName          = "myArtist"
	playlistName        = "My playlist"
	playlistDescription = "Some playlist"
	isPlaylistPublic    = true
)

func defaultContext() context.Context {
	return context.Background()
}

func defaultPlaylist() PlaylistDetails {
	return PlaylistDetails{Name: playlistName, Description: playlistDescription, IsPublic: isPlaylistPublic}
}

func defaultSongs() []interface{} {
	return []interface{}{
		song.NewSong("some_uri"),
		song.NewSong("another_uri"),
	}
}

func songsWithErrors() []interface{} {
	return []interface{}{
		errors.New("Some error"),
		song.NewSong("another_uri"),
	}
}

func errorSongs() []interface{} {
	return []interface{}{
		errors.New("Some error"),
		errors.New("Some other error"),
	}
}

func defaultSetlist() setlist.Setlist {
	songs := []setlist.Song{
		setlist.NewSong("My song"),
		setlist.NewSong("My other song"),
	}
	return setlist.NewSetlist(artistName, songs)
}

func emptySetlist() setlist.Setlist {
	return setlist.NewSetlist(artistName, []setlist.Song{})
}

func defaultGetSongArgs() []song.GetSongArgs {
	return []song.GetSongArgs{
		{Context: defaultContext(), Artist: artistName, Title: "My song"},
		{Context: defaultContext(), Artist: artistName, Title: "My other song"},
	}
}

func defaultAddSongsArgs() AddSongsArgs {
	return AddSongsArgs{
		Context:    defaultContext(),
		PlaylistId: playlistId,
		Songs: []song.Song{
			song.NewSong("some_uri"),
			song.NewSong("another_uri"),
		},
	}
}

func addSongsArgsWithErrors() AddSongsArgs {
	return AddSongsArgs{
		Context:    defaultContext(),
		PlaylistId: playlistId,
		Songs: []song.Song{
			song.NewSong("another_uri"),
		},
	}
}

func newFakeSetlistRepository() setlist.FakeSetlistRepository {
	repository := setlist.NewFakeSetlistRepository()
	repository.SetReturnValue(defaultSetlist())
	return repository
}

func newFakeSongRepository() song.FakeSongRepository {
	repository := song.NewFakeSongRepository()
	repository.SetSongs(defaultSongs())
	return repository
}

func testSetup() (FakePlaylistRepository, setlist.FakeSetlistRepository, song.FakeSongRepository) {
	playlistRepository := NewFakePlaylistRepository()
	playlistRepository.SetCreatedPlaylistId("some_id")
	setlistRepository := newFakeSetlistRepository()
	songRepository := newFakeSongRepository()
	return playlistRepository, setlistRepository, songRepository
}

func TestCreatePlaylistRepositoryCalledWithArgs(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	_, err := service.CreatePlaylist(defaultContext(), defaultPlaylist())

	actual := playlistRepository.GetCreatePlaylistArgs()
	expected := CreatePlaylistArgs{Context: defaultContext(), Playlist: defaultPlaylist()}
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestCreatePlaylistReturnsPlaylistIdOnSuccess(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	expected := "some random id"
	playlistRepository.SetCreatedPlaylistId(expected)
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	actual, err := service.CreatePlaylist(defaultContext(), defaultPlaylist())

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestCreatePlaylistReturnsErrorIfRepositoryErrors(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	playlistRepository.SetError(errors.New("test error"))
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	_, err := service.CreatePlaylist(defaultContext(), defaultPlaylist())

	assert.NotNil(t, err)
}

func TestAddSetlistSetlistRepositoryCalledWithArgs(t *testing.T) {
	artist := artistName
	minSongs := 6
	playlistRepository, setlistRepository, songRepository := testSetup()

	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)
	service.SetMinSongs(minSongs)

	err := service.AddSetlist(defaultContext(), playlistId, artist)

	actual := setlistRepository.GetGetSetlistArgs()
	expected := setlist.GetSetlistArgs{Artist: artist, MinSongs: minSongs}
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddSetlistReturnsErrorOnSetlistRepositoryError(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	returnError := errors.New("test error")
	setlistRepository.SetError(returnError)
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultContext(), playlistId, artistName)

	assert.NotNil(t, err)
}

func TestAddSetlistSongRepositoryCalledWithSetlistSongs(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultContext(), playlistId, artistName)

	actual := songRepository.GetGetSongArgs()
	expected := defaultGetSongArgs()
	assert.Nil(t, err)
	if !testtools.HaveSameElements(expected, actual) {
		t.Errorf("Expected called songs %v, found %v", expected, actual)
	}
}

func TestAddSetlistAddsSongsFetched(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(defaultSongs())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultContext(), playlistId, artistName)

	actual := playlistRepository.GetAddSongArgs()
	expected := defaultAddSongsArgs()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddSetlistAddsOnlySongsFetchedWithoutError(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(songsWithErrors())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultContext(), "myPlaylist", artistName)

	actual := playlistRepository.GetAddSongArgs()
	expected := addSongsArgsWithErrors()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddSetlistSetlistRaisesErrorIfSetlistEmpty(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(errorSongs())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultContext(), playlistId, artistName)

	assert.NotNil(t, err)
}

func TestAddSetlistSetlistRaisesErrorIfNoSongsFound(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	setlistRepository.SetReturnValue(emptySetlist())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultContext(), playlistId, artistName)

	assert.NotNil(t, err)
}
