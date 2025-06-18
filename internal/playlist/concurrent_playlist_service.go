package playlist

import (
	"context"
	"errors"

	"festwrap/internal/logging"
	"festwrap/internal/setlist"
	"festwrap/internal/song"
	"fmt"
	"time"
)

type FetchSongResult struct {
	Song song.Song
	Err  error
}

type ConcurrentPlaylistService struct {
	playlistRepository PlaylistRepository
	setlistRepository  setlist.SetlistRepository
	songRepository     song.SongRepository
	minSongs           int
	addSetlistSleepMs  int
	logger             logging.Logger
}

func NewConcurrentPlaylistService(
	playlistRepository PlaylistRepository,
	setlistRepository setlist.SetlistRepository,
	songRepository song.SongRepository,
	logger logging.Logger,
) ConcurrentPlaylistService {
	return ConcurrentPlaylistService{
		playlistRepository: playlistRepository,
		setlistRepository:  setlistRepository,
		songRepository:     songRepository,
		logger:             logger,
		minSongs:           4,
		addSetlistSleepMs:  0,
	}
}

func (s *ConcurrentPlaylistService) CreatePlaylistWithArtists(
	ctx context.Context,
	playlist PlaylistDetails,
	artists []string,
) (PlaylistCreation, error) {
	playlistId, err := s.playlistRepository.CreatePlaylist(ctx, playlist)
	if err != nil {
		return PlaylistCreation{}, errors.New("could not create playlist")
	}

	errors := 0
	for i, artist := range artists {
		if i > 0 {
			// Sleep to avoid hitting Setlistfm rate limit
			time.Sleep(time.Duration(s.addSetlistSleepMs) * time.Millisecond)
		}

		err := s.addSetlistToPlaylist(ctx, playlistId, artist)
		if err != nil {
			message := fmt.Sprintf("could not add songs for %s to playlist %s: %v", artist, playlistId, err)
			s.logger.Warn(message)
			errors += 1
		}
	}
	if errors == len(artists) {
		s.logger.Error(fmt.Sprintf("could not add any of artists %v to playlist %s", artists, playlistId))
		return PlaylistCreation{}, fmt.Errorf("all artists failed to be added to playlist %s", playlistId)
	}

	var status CreationStatus
	if errors == 0 {
		status = Success
	} else {
		status = PartialFailure
	}
	return PlaylistCreation{PlaylistId: playlistId, Status: status}, nil
}

func (s *ConcurrentPlaylistService) addSetlistToPlaylist(ctx context.Context, playlistId string, artist string) error {
	setlist, err := s.setlistRepository.GetSetlist(artist, s.minSongs)
	if err != nil {
		return err
	}

	ch := make(chan FetchSongResult)
	for _, song := range setlist.GetSongs() {
		go s.fetchSong(ctx, artist, song, ch)
	}

	songs := []song.Song{}
	for i := 0; i < len(setlist.GetSongs()); i++ {
		result := <-ch
		if result.Err == nil {
			songs = append(songs, result.Song)
		}
	}

	if len(songs) == 0 {
		return fmt.Errorf("no songs to add to playlist %s for artist %s", playlistId, artist)
	}

	err = s.playlistRepository.AddSongs(ctx, playlistId, songs)
	if err != nil {
		return err
	}

	return nil
}

func (s *ConcurrentPlaylistService) SetMinSongs(minSongs int) {
	s.minSongs = minSongs
}

func (s *ConcurrentPlaylistService) fetchSong(
	ctx context.Context,
	artist string,
	song setlist.Song,
	ch chan<- FetchSongResult,
) {
	songDetails, err := s.songRepository.GetSong(ctx, artist, song.GetTitle())
	ch <- FetchSongResult{Song: songDetails, Err: err}
}
