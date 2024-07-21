package setlist

type Setlist struct {
	artist string
	songs  []Song
}

func (s Setlist) GetArtist() string {
	return s.artist
}

func (s Setlist) GetSongs() []Song {
	return s.songs
}

func NewSetlist(artist string, songs []Song) Setlist {
	return Setlist{artist: artist, songs: songs}
}
