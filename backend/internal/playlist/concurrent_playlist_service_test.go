package playlist

import (
	"errors"
	"reflect"
	"testing"

	"festwrap/internal/setlist"
	"festwrap/internal/song"
	"festwrap/internal/testtools"
)

func PlaylistId() string {
	return "myPlaylist"
}

func Artist() string {
	return "myArtist"
}

func Songs() []interface{} {
	return []interface{}{
		song.NewSong("some_uri"),
		song.NewSong("another_uri"),
	}
}

func SongsWithErrors() []interface{} {
	return []interface{}{
		errors.New("Some error"),
		song.NewSong("another_uri"),
	}
}

func ErrorSongs() []interface{} {
	return []interface{}{
		errors.New("Some error"),
		errors.New("Some other error"),
	}
}

func Setlist() setlist.Setlist {
	songs := []setlist.Song{
		setlist.NewSong("My song"),
		setlist.NewSong("My other song"),
	}
	return setlist.NewSetlist(Artist(), songs)
}

func EmptySetlist() setlist.Setlist {
	return setlist.NewSetlist(Artist(), []setlist.Song{})
}

func DefaultGetSongArgs() []song.GetSongArgs {
	return []song.GetSongArgs{
		song.GetSongArgs{Artist: Artist(), Title: "My song"},
		song.GetSongArgs{Artist: Artist(), Title: "My other song"},
	}
}

func GetSongArgsWithErrors() []song.GetSongArgs {
	return []song.GetSongArgs{
		song.GetSongArgs{Artist: Artist(), Title: "My other song"},
	}
}

func DefaultAddSongsArgs() AddSongsArgs {
	return AddSongsArgs{
		PlaylistId: PlaylistId(),
		Songs: []song.Song{
			song.NewSong("some_uri"),
			song.NewSong("another_uri"),
		},
	}
}

func AddSongsArgsWithErrors() AddSongsArgs {
	return AddSongsArgs{
		PlaylistId: PlaylistId(),
		Songs: []song.Song{
			song.NewSong("another_uri"),
		},
	}
}

func NewFakeSetlistRepository() setlist.FakeSetlistRepository {
	repository := setlist.NewFakeSetlistRepository()
	returnSetlist := Setlist()
	repository.SetReturnValue(&returnSetlist)
	return repository
}

func NewFakeSongRepository() song.FakeSongRepository {
	repository := song.NewFakeSongRepository()
	repository.SetSongs(Songs())
	return repository
}

func testSetup() (FakePlaylistRepository, setlist.FakeSetlistRepository, song.FakeSongRepository) {
	playlistRepository := NewFakePlaylistRepository()
	setlistRepository := NewFakeSetlistRepository()
	songRepository := NewFakeSongRepository()
	return playlistRepository, setlistRepository, songRepository
}

func TestAddSetlistSetlistRepositoryCalledWithArtist(t *testing.T) {
	expected := Artist()
	playlistRepository, setlistRepository, songRepository := testSetup()

	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(PlaylistId(), expected)

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

	err := service.AddSetlist(PlaylistId(), Artist())

	if err != returnError {
		t.Errorf("Setlist repository should have returned an error but it did not")
	}
}

func TestAddSetlistSongRepositoryCalledWithSetlistSongs(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(PlaylistId(), Artist())

	actual := songRepository.GetGetSongArgs()
	expected := DefaultGetSongArgs()
	testtools.AssertErrorIsNil(t, err)
	if !testtools.HaveSameElements[song.GetSongArgs](actual, expected) {
		t.Errorf("Expected called songs %v, found %v", expected, actual)
	}
}

func TestAddSetlistAddsSongsFetched(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(Songs())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(PlaylistId(), Artist())

	actual := playlistRepository.GetAddSongArgs()
	expected := DefaultAddSongsArgs()
	testtools.AssertErrorIsNil(t, err)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected added songs call to be %v, found %v", expected, actual)
	}
}

func TestAddSetlistAddsOnlySongsFetchedWithoutError(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(SongsWithErrors())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist("myPlaylist", Artist())

	actual := playlistRepository.GetAddSongArgs()
	expected := AddSongsArgsWithErrors()
	testtools.AssertErrorIsNil(t, err)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected added songs call to be %v, found %v", expected, actual)
	}
}

func TestAddSetlistSetlistRaisesErrorIfSetlistEmpty(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	songRepository.SetSongs(ErrorSongs())
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(PlaylistId(), Artist())

	testtools.AssertErrorNotNil(t, err)
}

func TestAddSetlistSetlistRaisesErrorIfNoSongsFound(t *testing.T) {
	playlistRepository, setlistRepository, songRepository := testSetup()
	setlist := EmptySetlist()
	setlistRepository.SetReturnValue(&setlist)
	service := NewConcurrentPlaylistService(&playlistRepository, &setlistRepository, &songRepository)

	err := service.AddSetlist(PlaylistId(), Artist())

	testtools.AssertErrorNotNil(t, err)
}
