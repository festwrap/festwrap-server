package song

import "context"

type SongRepository interface {
	GetSong(ctx context.Context, artist string, title string) (*Song, error)
}
