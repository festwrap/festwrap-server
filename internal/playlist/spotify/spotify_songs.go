package spotify

import (
	"festwrap/internal/song"
)

type SpotifySongs struct {
	Uris []string `json:"uris"`
}

func NewSpotifySongs(songs []song.Song) SpotifySongs {
	songUris := []string{}
	for _, currentSong := range songs {
		songUris = append(songUris, currentSong.GetUri())
	}
	return SpotifySongs{Uris: songUris}
}
