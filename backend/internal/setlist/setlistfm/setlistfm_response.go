package setlistfm

import (
	"festwrap/internal/setlist"
)

type setlistfmSong struct {
	Name string `json:"name"`
}

type setlistfmSet struct {
	Songs []setlistfmSong `json:"song"`
}

type setlistfmArtist struct {
	Name string `json:"name"`
}

type setlistFMSets struct {
	Sets []setlistfmSet `json:"set"`
}

type setlistFMSetlist struct {
	Artist setlistfmArtist `json:"artist"`
	Sets   setlistFMSets   `json:"sets"`
}

func (s *setlistFMSetlist) GetSongs() []setlist.Song {
	songs := []setlist.Song{}
	for _, set := range s.Sets.Sets {
		for _, song := range set.Songs {
			songs = append(songs, setlist.NewSong(song.Name))
		}
	}
	return songs
}

type SetlistFMResponse struct {
	Body []setlistFMSetlist `json:"setlist"`
}
