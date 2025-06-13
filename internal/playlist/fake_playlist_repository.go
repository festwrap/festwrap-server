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
	Playlist Playlist
}

type FakePlaylistRepository struct {
	addSongArgs        AddSongsArgs
	createPlaylistArgs CreatePlaylistArgs
	searchedPlaylists  []Playlist
	createdPlaylistId  string
	err                error
}

func NewFakePlaylistRepository() FakePlaylistRepository {
	return FakePlaylistRepository{searchedPlaylists: []Playlist{}}
}

func (s *FakePlaylistRepository) CreatePlaylist(ctx context.Context, playlist Playlist) (string, error) {
	s.createPlaylistArgs = CreatePlaylistArgs{Context: ctx, Playlist: playlist}
	return s.createdPlaylistId, s.err
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

func (s *FakePlaylistRepository) GetCreatePlaylistArgs() CreatePlaylistArgs {
	return s.createPlaylistArgs
}

func (s *FakePlaylistRepository) SetSearchedPlaylists(playlists []Playlist) {
	s.searchedPlaylists = playlists
}

func (s *FakePlaylistRepository) SetCreatedPlaylistId(id string) {
	s.createdPlaylistId = id
}
