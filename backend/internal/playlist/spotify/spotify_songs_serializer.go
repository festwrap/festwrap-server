package spotify

import (
	"encoding/json"
	"festwrap/internal/song"
)

type SpotifySongsSerializer struct{}

func (s *SpotifySongsSerializer) Serialize(songs []song.Song) ([]byte, error) {
	songUris := []string{}
	for _, currentSong := range songs {
		songUris = append(songUris, currentSong.GetUri())
	}
	return json.Marshal(SpotifyTrackUris{Uris: songUris})
}
