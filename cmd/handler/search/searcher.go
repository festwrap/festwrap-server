package search

import "context"

type Searcher[T any] interface {
	Search(ctx context.Context, name string, limit int) ([]T, error)
}
