package song

type Song struct {
	uri string
}

func (s *Song) GetUri() string {
	return s.uri
}

func NewSong(uri string) Song {
	return Song{uri: uri}
}
