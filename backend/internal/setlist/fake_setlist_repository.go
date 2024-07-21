package setlist

type FakeSetlistRepository struct {
	artist      string
	returnValue *Setlist
	err         error
}

func (s *FakeSetlistRepository) GetSetlist(artist string) (*Setlist, error) {
	s.artist = artist
	return s.returnValue, s.err
}

func (s *FakeSetlistRepository) GetCalledArtist() string {
	return s.artist
}

func (s *FakeSetlistRepository) SetReturnValue(setlist *Setlist) {
	s.returnValue = setlist
}

func (s *FakeSetlistRepository) SetError(err error) {
	s.err = err
}

func NewFakeSetlistRepository() FakeSetlistRepository {
	return FakeSetlistRepository{}
}
