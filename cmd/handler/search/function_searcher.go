package search

import "context"

type FunctionSearcher[T any] struct {
	fn func(ctx context.Context, name string, limit int) ([]T, error)
}

func NewFunctionSearcher[T any](fn func(ctx context.Context, name string, limit int) ([]T, error)) FunctionSearcher[T] {
	return FunctionSearcher[T]{fn: fn}
}

func (s *FunctionSearcher[T]) Search(ctx context.Context, name string, limit int) ([]T, error) {
	return s.fn(ctx, name, limit)
}
