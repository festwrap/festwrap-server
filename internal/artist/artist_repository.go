package artist

import "context"

type ArtistRepository interface {
	SearchArtist(ctx context.Context, name string, limit int) ([]Artist, error)
}
