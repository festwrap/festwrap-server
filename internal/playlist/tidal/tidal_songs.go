package tidal

import (
	"festwrap/internal/song"
)

type tidalSong struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type tidalSongs struct {
	Data []tidalSong `json:"data"`
}

func NewTidalSongs(songs []song.Song) tidalSongs {
	result := []tidalSong{}
	for _, currentSong := range songs {
		result = append(result, tidalSong{Id: currentSong.GetUri(), Type: "tracks"})
	}
	return tidalSongs{Data: result}
}
