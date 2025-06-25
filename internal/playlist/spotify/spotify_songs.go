package spotify

import (
	"festwrap/internal/song"
)

type spotifySongs struct {
	Uris []string `json:"uris"`
}

func NewSpotifySongs(songs []song.Song) spotifySongs {
	songUris := []string{}
	for _, currentSong := range songs {
		songUris = append(songUris, currentSong.GetUri())
	}
	return spotifySongs{Uris: songUris}
}
