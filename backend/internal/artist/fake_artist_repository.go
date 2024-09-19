package artist

import "context"

type SearchArtistArgs struct {
	Context context.Context
	Name    string
	Limit   int
}

type SearchArtistReturn struct {
	value []Artist
	err   error
}

type FakeArtistRepository struct {
	searchArgs   SearchArtistArgs
	searchReturn SearchArtistReturn
}

func (r *FakeArtistRepository) SearchArtist(ctx context.Context, name string, limit int) ([]Artist, error) {
	r.searchArgs = SearchArtistArgs{Name: name, Limit: limit, Context: ctx}
	return r.searchReturn.value, r.searchReturn.err
}

func (r *FakeArtistRepository) GetSearchArtistArgs() SearchArtistArgs {
	return r.searchArgs
}

func (r *FakeArtistRepository) SetSearchReturnValue(value []Artist) {
	r.searchReturn.value = value
}

func (r *FakeArtistRepository) SetSearchArtistError(err error) {
	r.searchReturn.err = err
}
