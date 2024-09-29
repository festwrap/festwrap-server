package setlist

type FakeSetlistRepository struct {
	getArgs  GetSetlistArgs
	getValue getSetlistValue
}

type GetSetlistArgs struct {
	Artist   string
	MinSongs int
}

type getSetlistValue struct {
	response Setlist
	err      error
}

func NewFakeSetlistRepository() FakeSetlistRepository {
	return FakeSetlistRepository{}
}

func (s *FakeSetlistRepository) GetSetlist(artist string, minSongs int) (*Setlist, error) {
	s.getArgs = GetSetlistArgs{Artist: artist, MinSongs: minSongs}
	return &s.getValue.response, s.getValue.err
}

func (s *FakeSetlistRepository) GetGetSetlistArgs() GetSetlistArgs {
	return s.getArgs
}

func (s *FakeSetlistRepository) SetReturnValue(setlist Setlist) {
	s.getValue.response = setlist
}

func (s *FakeSetlistRepository) SetError(err error) {
	s.getValue.err = err
}
