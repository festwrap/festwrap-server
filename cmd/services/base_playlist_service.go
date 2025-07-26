package playlist

import (
	"context"

	"festwrap/internal/event"
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
}

type BasePlaylistService struct {
	playlistRepository       playlist.PlaylistRepository
	setlistRepository        setlist.SetlistRepository
	songRepository           song.SongRepository
	playlistCreationNotifier event.Notifier[event.PlaylistCreatedEvent]
	minSongs                 int
	addSetlistSleepMs        int
	logger                   logging.Logger
}

func NewBasePlaylistService(
	playlistRepository playlist.PlaylistRepository,
	setlistRepository setlist.SetlistRepository,
	songRepository song.SongRepository,
	logger logging.Logger,
) BasePlaylistService {
	return BasePlaylistService{
		playlistRepository:       playlistRepository,
		setlistRepository:        setlistRepository,
		songRepository:           songRepository,
		playlistCreationNotifier: event.NewBaseNotifier[event.PlaylistCreatedEvent](),
		logger:                   logger,
		minSongs:                 4,
		addSetlistSleepMs:        0,
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

	s.notifyPlaylistCreated(ctx, playlistId, playlist.Name, artists, status)

	return PlaylistCreation{PlaylistId: playlistId, Status: status}, nil
}

func (s *BasePlaylistService) SetAddSetlistSleep(sleepMs int) {
	s.addSetlistSleepMs = sleepMs
}

func (s *BasePlaylistService) SetMinSongs(minSongs int) {
	s.minSongs = minSongs
}

func (s *BasePlaylistService) SetPlaylistCreateNotifier(
	subject event.Notifier[event.PlaylistCreatedEvent],
) *BasePlaylistService {
	s.playlistCreationNotifier = subject
	return s
}

func (s *BasePlaylistService) addSetlistToPlaylist(ctx context.Context, playlistId string, artist string) error {
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

func (s *BasePlaylistService) fetchSong(
	ctx context.Context,
	artist string,
	song setlist.Song,
	ch chan<- FetchSongResult,
) {
	songDetails, err := s.songRepository.GetSong(ctx, artist, song.GetTitle())
	ch <- FetchSongResult{Song: songDetails, Err: err}
}

func (s *BasePlaylistService) notifyPlaylistCreated(
	ctx context.Context,
	playlistId,
	playlistName string,
	artists []string,
	status CreationStatus,
) {
	var eventStatus event.PlaylistCreationStatus
	if status == Success {
		eventStatus = event.PLAYLIST_CREATED_OK
	} else if status == PartialFailure {
		eventStatus = event.PLAYLIST_CREATED_PARTIAL_ERRORS
	} else {
		eventStatus = event.PLAYLIST_CREATED_OK
		s.logger.Error(fmt.Sprintf("unknown playlist creation status: %v. Using ok status as default", status))
	}

	playlistCreatedEvent := event.PlaylistCreatedEvent{
		Playlist: event.CreatedPlaylist{
			Id:      playlistId,
			Name:    playlistName,
			Artists: artists,
			Type:    event.PLAYLIST_TYPE_SPOTIFY,
		},
		CreationStatus: eventStatus,
	}
	s.playlistCreationNotifier.Notify(event.NewEventWrapper(playlistCreatedEvent))
}
