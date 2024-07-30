package playlist

import (
	"errors"
	"reflect"
	"testing"

	"festwrap/internal/setlist"
	"festwrap/internal/song"
	"festwrap/internal/testtools"
)

func defaultPlaylistId() string {
	return "myPlaylist"
}

func defaultArtist() string {
	return "myArtist"
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
	return setlist.NewSetlist(defaultArtist(), songs)
}

func emptySetlist() setlist.Setlist {
	return setlist.NewSetlist(defaultArtist(), []setlist.Song{})
}

func defaultGetSongArgs() []song.GetSongArgs {
	return []song.GetSongArgs{
		{Artist: defaultArtist(), Title: "My song"},
		{Artist: defaultArtist(), Title: "My other song"},
	}
}

func defaultAddSongsArgs() AddSongsArgs {
	return AddSongsArgs{
		PlaylistId: defaultPlaylistId(),
		Songs: []song.Song{
			song.NewSong("some_uri"),
			song.NewSong("another_uri"),
		},
	}
}

func addSongsArgsWithErrors() AddSongsArgs {
	return AddSongsArgs{
		PlaylistId: defaultPlaylistId(),
		Songs: []song.Song{
			song.NewSong("another_uri"),
		},
	}
}

func newFakeSetlistRepository() setlist.FakeSetlistRepository {
	repository := setlist.NewFakeSetlistRepository()
	returnSetlist := defaultSetlist()
	repository.SetReturnValue(&returnSetlist)
	return repository
}

func newFakeSongRepository() song.FakeSongRepository {
	repository := song.NewFakeSongRepository()
	repository.SetSongs(defaultSongs())
	return repository
}

func testSetup() (FakePlaylistRepository, setlist.FakeSetlistRepository, song.FakeSongRepository) {
	playlistRepository := NewFakePlaylistRepository()
	setlistRepository := newFakeSetlistRepository()
	songRepository := newFakeSongRepository()
	return playlistRepository, setlistRepository, songRepository
}

func TestAddSetlistSetlistRepositoryCalledWithArtist(t *testing.T) {
	expected := defaultArtist()
	playlistRepository, setlistRepository, songRepository := testSetup()

	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultPlaylistId(), expected)

	actual := setlistRepository.GetCalledArtist()
	testtools.AssertErrorIsNil(t, err)
	if actual != expected {
		t.Errorf("Setlist repository to be called with %s, found %s", expected, actual)
	}
}

func TestAddSetlistReturnsErrorOnSetlistRepositoryError(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	returnError := errors.New("test error")
	setlistRepository.SetError(returnError)
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultPlaylistId(), defaultArtist())

	if err != returnError {
		t.Errorf("Setlist repository should have returned an error but it did not")
	}
}

func TestAddSetlistSongRepositoryCalledWithSetlistSongs(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultPlaylistId(), defaultArtist())

	actual := songRepository.GetGetSongArgs()
	expected := defaultGetSongArgs()
	testtools.AssertErrorIsNil(t, err)
	if !testtools.HaveSameElements[song.GetSongArgs](actual, expected) {
		t.Errorf("Expected called songs %v, found %v", expected, actual)
	}
}

func TestAddSetlistAddsSongsFetched(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(defaultSongs())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultPlaylistId(), defaultArtist())

	actual := playlistRepository.GetAddSongArgs()
	expected := defaultAddSongsArgs()
	testtools.AssertErrorIsNil(t, err)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected added songs call to be %v, found %v", expected, actual)
	}
}

func TestAddSetlistAddsOnlySongsFetchedWithoutError(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(songsWithErrors())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist("myPlaylist", defaultArtist())

	actual := playlistRepository.GetAddSongArgs()
	expected := addSongsArgsWithErrors()
	testtools.AssertErrorIsNil(t, err)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected added songs call to be %v, found %v", expected, actual)
	}
}

func TestAddSetlistSetlistRaisesErrorIfSetlistEmpty(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(errorSongs())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultPlaylistId(), defaultArtist())

	testtools.AssertErrorNotNil(t, err)
}

func TestAddSetlistSetlistRaisesErrorIfNoSongsFound(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	setlist := emptySetlist()
	setlistRepository.SetReturnValue(&setlist)
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(defaultPlaylistId(), defaultArtist())

	testtools.AssertErrorNotNil(t, err)
}
