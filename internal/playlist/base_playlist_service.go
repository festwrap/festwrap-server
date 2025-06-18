package playlist

import (
	"context"
	"festwrap/internal/logging"
	"festwrap/internal/setlist"
	"festwrap/internal/song"
	"fmt"
	"time"
)

type FetchSongResult struct {
	Song *song.Song
	Err  error
}

type BasePlaylistService struct {
	playlistRepository PlaylistRepository
	setlistRepository  setlist.SetlistRepository
	songRepository     song.SongRepository
	minSongs           int
	addSetlistSleepMs  int
	logger             logging.Logger
}

func NewBasePlaylistService(
	playlistRepository PlaylistRepository,
	setlistRepository setlist.SetlistRepository,
	songRepository song.SongRepository,
	addSetlistSleepMs int,
	logger logging.Logger,
) BasePlaylistService {
	return BasePlaylistService{
		playlistRepository: playlistRepository,
		setlistRepository:  setlistRepository,
		songRepository:     songRepository,
		minSongs:           4,
		logger:             logger,
		addSetlistSleepMs:  addSetlistSleepMs,
	}
}

func (s *BasePlaylistService) CreatePlaylistWithArtists(ctx context.Context, playlist Playlist, artists []PlaylistArtist) (PlaylistCreation, error) {
	playlistId, err := s.playlistRepository.CreatePlaylist(ctx, playlist)
	if err != nil {
		return PlaylistCreation{}, fmt.Errorf("could not create playlist: %v", err)
	}

	errorCount := 0
	for i, artist := range artists {
		if i > 0 {
			// Sleep to avoid hitting Setlistfm rate limit
			time.Sleep(time.Duration(s.addSetlistSleepMs) * time.Millisecond)
		}

		err := s.addSetlist(ctx, playlistId, artist.Name)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("could not add setlist for artist %s to playlist %s: %v", artist.Name, playlistId, err))
			errorCount += 1
		}
	}

	if errorCount == 0 {
		return PlaylistCreation{PlaylistId: playlistId, Status: Success}, nil
	} else if errorCount > 0 && errorCount < len(artists) {
		return PlaylistCreation{PlaylistId: playlistId, Status: PartialFailure}, nil
	} else {
		return PlaylistCreation{}, fmt.Errorf("could not add any songs to playlist %s", playlistId)
	}
}

func (s *BasePlaylistService) addSetlist(ctx context.Context, playlistId string, artist string) error {
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
			songs = append(songs, *result.Song)
		}
	}

	if len(songs) == 0 {
		return fmt.Errorf("no songs found for artist %s", artist)
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
	ch chan<- FetchSongResult,
) {
	songDetails, err := s.songRepository.GetSong(ctx, artist, song.GetTitle())
	ch <- FetchSongResult{Song: songDetails, Err: err}
}
