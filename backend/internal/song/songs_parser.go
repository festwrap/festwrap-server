package song

type SongsParser interface {
	Parse(songs []byte) (*[]Song, error)
}
