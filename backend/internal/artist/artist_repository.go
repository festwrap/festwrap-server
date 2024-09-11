package artist

type ArtistRepository interface {
	SearchArtist(name string, limit int) (*[]Artist, error)
}
