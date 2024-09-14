package song

type Song struct {
	uri string
}

func NewSong(uri string) Song {
	return Song{uri: uri}
}

func (s *Song) GetUri() string {
	return s.uri
}
