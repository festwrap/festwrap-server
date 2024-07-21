package song

import (
	"fmt"
	"sync"
)

type GetSongArgs struct {
	Artist string
	Title  string
}

type FakeSongRepository struct {
	// Creating a wrapped structure so we do not need
	// to force interface to be referenced as a pointer
	repository *WrappedFakeSongRepository
}

func (r *FakeSongRepository) GetSong(artist string, title string) (*Song, error) {
	return r.repository.GetSong(artist, title)
}

func (r *FakeSongRepository) GetGetSongArgs() []GetSongArgs {
	return r.repository.getSongArgs
}

func (r *FakeSongRepository) SetSongs(songs []interface{}) {
	r.repository.songs = songs
}

func NewFakeSongRepository() FakeSongRepository {
	return FakeSongRepository{
		repository: &WrappedFakeSongRepository{getSongArgs: []GetSongArgs{}, songs: []interface{}{}},
	}
}

type WrappedFakeSongRepository struct {
	getSongArgs []GetSongArgs
	songs       []interface{}
	mutex       sync.Mutex
}

func (w *WrappedFakeSongRepository) GetSong(artist string, title string) (*Song, error) {
	w.getSongArgs = append(w.getSongArgs, GetSongArgs{Artist: artist, Title: title})
	return w.popSongLeft()
}

func (w *WrappedFakeSongRepository) popSongLeft() (*Song, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(w.songs) == 0 {
		panic("Fake repository has not songs left")
	}

	top := w.songs[0]
	w.songs = w.songs[1:]

	switch result := top.(type) {
	case Song:
		return &result, nil
	case error:
		return nil, result
	default:
		message := fmt.Sprintf("Fake repository should only return errors or songs. Found %v", top)
		panic(message)
	}
}
