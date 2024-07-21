package song

type SongRepository interface {
	GetSong(artist string, title string) (*Song, error)
}
