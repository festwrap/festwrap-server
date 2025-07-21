package playlist

import (
	"context"

	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	"festwrap/internal/setlist"
	"festwrap/internal/song"
	"fmt"
	"time"
)

type FetchSongResult struct {
	Song song.Song
	Err  error
	Rank int
}

type BasePlaylistService struct {
	playlistRepository playlist.PlaylistRepository
	setlistRepository  setlist.SetlistRepository
	songRepository     song.SongRepository
	minSongs           int
	addSetlistSleepMs  int
	logger             logging.Logger
}

func NewBasePlaylistService(
	playlistRepository playlist.PlaylistRepository,
	setlistRepository setlist.SetlistRepository,
	songRepository song.SongRepository,
	logger logging.Logger,
) BasePlaylistService {
	return BasePlaylistService{
		playlistRepository: playlistRepository,
		setlistRepository:  setlistRepository,
		songRepository:     songRepository,
		logger:             logger,
		minSongs:           4,
		addSetlistSleepMs:  0,
	}
}

func (s *BasePlaylistService) CreatePlaylistWithArtists(
	ctx context.Context,
	playlist playlist.PlaylistDetails,
	artists []string,
) (PlaylistCreation, error) {
	playlistId, err := s.playlistRepository.CreatePlaylist(ctx, playlist)
	if err != nil {
		return PlaylistCreation{}, fmt.Errorf("could not create playlist: %v", err)
	}

	errors := 0
	for i, artist := range artists {
		if i > 0 {
			// Sleep to avoid hitting Setlistfm rate limit
			time.Sleep(time.Duration(s.addSetlistSleepMs) * time.Millisecond)
		}
		err := s.addSetlistToPlaylist(ctx, playlistId, artist)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("could not add songs for %s to playlist %s: %v", artist, playlistId, err))
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

func (s *BasePlaylistService) SetAddSetlistSleep(sleepMs int) {
	s.addSetlistSleepMs = sleepMs
}

func (s *BasePlaylistService) addSetlistToPlaylist(ctx context.Context, playlistId string, artist string) error {
	setlist, err := s.setlistRepository.GetSetlist(artist, s.minSongs)
	if err != nil {
		return err
	}

	songsCount := len(setlist.GetSongs())
	ch := make(chan FetchSongResult)
	rankedResults := make([]FetchSongResult, songsCount)
	for i, song := range setlist.GetSongs() {
		go s.fetchSong(ctx, artist, song, i, ch)
	}

	// Keep songs in the original setlist order
	for range songsCount {
		fetchResult := <-ch
		rankedResults[fetchResult.Rank] = fetchResult
	}

	songs := []song.Song{}
	for _, fetchResult := range rankedResults {
		if fetchResult.Err == nil {
			songs = append(songs, fetchResult.Song)
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

func (s *BasePlaylistService) SetMinSongs(minSongs int) {
	s.minSongs = minSongs
}

func (s *BasePlaylistService) fetchSong(
	ctx context.Context,
	artist string,
	song setlist.Song,
	rank int,
	ch chan<- FetchSongResult,
) {
	songDetails, err := s.songRepository.GetSong(ctx, artist, song.GetTitle())
	ch <- FetchSongResult{Song: songDetails, Err: err, Rank: rank}
}
