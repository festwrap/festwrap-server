package playlist

import (
	"context"
	"festwrap/internal/song"
)

type AddSongsArgs struct {
	Context    context.Context
	PlaylistId string
	Songs      []song.Song
}

type CreatePlaylistArgs struct {
	Context  context.Context
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

func (s *FakePlaylistRepository) CreatePlaylist(ctx context.Context, userId string, playlist Playlist) error {
	s.createPlaylistArgs = CreatePlaylistArgs{Context: ctx, UserId: userId, Playlist: playlist}
	return s.err
}

func (s *FakePlaylistRepository) AddSongs(ctx context.Context, playlistId string, songs []song.Song) error {
	s.addSongArgs = AddSongsArgs{Context: ctx, PlaylistId: playlistId, Songs: songs}
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
