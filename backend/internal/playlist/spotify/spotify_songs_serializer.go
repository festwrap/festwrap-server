package spotify

import (
	"encoding/json"
	"festwrap/internal/song"
)

type SongList struct {
	songs []song.Song
}

type SpotifyTrackUris struct {
	Uris []string `json:"uris"`
}

type SpotifySongsSerializer struct{}

func (s *SpotifySongsSerializer) Serialize(songList SongList) ([]byte, error) {
	songUris := []string{}
	for _, currentSong := range songList.songs {
		songUris = append(songUris, currentSong.GetUri())
	}
	return json.Marshal(SpotifyTrackUris{Uris: songUris})
}
