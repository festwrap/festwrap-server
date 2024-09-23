package playlist

import "festwrap/internal/song"

type AddSongsArgs struct {
	PlaylistId string
	Songs      []song.Song
}

type CreatePlaylistArgs struct {
	UserId   string
	Playlist Playlist
}

type FakePlaylistRepository struct {
	addSongArgs        AddSongsArgs
	createPlaylistArgs CreatePlaylistArgs
	err                error
}

func NewFakePlaylistRepository() FakePlaylistRepository {
	return FakePlaylistRepository{}
}

func (s *FakePlaylistRepository) CreatePlaylist(userId string, playlist Playlist) error {
	s.createPlaylistArgs = CreatePlaylistArgs{UserId: userId, Playlist: playlist}
	return s.err
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

func (s *FakePlaylistRepository) GetCreatePlaylistSongArgs() CreatePlaylistArgs {
	return s.createPlaylistArgs
}
