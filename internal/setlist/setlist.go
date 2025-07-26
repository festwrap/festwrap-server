package setlist

type Setlist struct {
	url    string
	artist string
	songs  []Song
}

func NewSetlist(artist string, songs []Song, url string) Setlist {
	return Setlist{artist: artist, songs: songs, url: url}
}

func (s Setlist) GetArtist() string {
	return s.artist
}

func (s *Setlist) SetArtist(artist string) {
	s.artist = artist
}

func (s Setlist) GetSongs() []Song {
	return s.songs
}

func (s *Setlist) GetUrl() string {
	return s.url
}
