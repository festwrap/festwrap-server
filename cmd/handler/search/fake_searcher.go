package search

import "context"

type SearchArgs struct {
	Context context.Context
	Name    string
	Limit   int
}

type SearchResult[T any] struct {
	Result []T
	Error  error
}

type FakeSearcher[T any] struct {
	searchArgs    SearchArgs
	searchResults SearchResult[T]
}

func NewFakeSearcher[T any]() *FakeSearcher[T] {
	return &FakeSearcher[T]{searchArgs: SearchArgs{}, searchResults: SearchResult[T]{}}
}

func (s *FakeSearcher[T]) Search(ctx context.Context, name string, limit int) ([]T, error) {
	s.searchArgs = SearchArgs{Context: ctx, Name: name, Limit: limit}
	return s.searchResults.Result, s.searchResults.Error
}

func (s *FakeSearcher[T]) SetSearchResult(value []T) {
	s.searchResults.Result = value
}

func (s *FakeSearcher[T]) SetSearchError(err error) {
	s.searchResults.Error = err
}

func (s *FakeSearcher[T]) GetSearchArgs() SearchArgs {
	return s.searchArgs
}
