package playlist

import "festwrap/internal/song"

type AddSongsArgs struct {
	PlaylistId string
	Songs      []song.Song
}

type FakePlaylistRepository struct {
	addSongArgs AddSongsArgs
	err         error
}

func (s *FakePlaylistRepository) AddSongs(playlistId string, songs []song.Song) error {
	s.addSongArgs = AddSongsArgs{PlaylistId: playlistId, Songs: songs}
	return s.err
}

func (s *FakePlaylistRepository) SetError(err error) {
	s.err = err
}

func (s *FakePlaylistRepository) GetAddSongArgs() AddSongsArgs {
	return s.addSongArgs
}

func NewFakePlaylistRepository() FakePlaylistRepository {
	return FakePlaylistRepository{}
}
