package search

import (
	"fmt"
	"net/http"
	"strconv"

	"festwrap/internal/logging"
	"festwrap/internal/serialization"
)

type SearchHandler[T any] struct {
	encoder      serialization.Encoder[[]T]
	defaultLimit int
	maxLimit     int
	searcher     Searcher[T]
	entityType   string
	logger       logging.Logger
}

func NewSearchHandler[T any](searcher Searcher[T], entityType string, logger logging.Logger) SearchHandler[T] {
	return SearchHandler[T]{
		encoder:      serialization.NewJsonEncoder[[]T](),
		searcher:     searcher,
		entityType:   entityType,
		maxLimit:     10,
		defaultLimit: 5,
		logger:       logger,
	}
}

func (h *SearchHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		message := fmt.Sprintf("Validation error: %s name was not provided", h.entityType)
		h.logger.Warn(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	limit, err := h.readLimit(r)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid limit for request: %v", err.Error()))
		http.Error(
			w,
			fmt.Sprintf("Validation error: invalid limit. It should be an integer in interval [1, %d]", h.maxLimit),
			http.StatusUnprocessableEntity,
		)
		return
	}
	h.logger.Info(fmt.Sprintf("Received new request for %s, using limit %d", name, limit))

	results, err := h.searcher.Search(r.Context(), name, limit)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error searching for %s: %v", h.entityType, err.Error()))
		h.writeUnexpectedError(w)
		return
	}
	h.logger.Info(fmt.Sprintf("Found %s %v for %s, using limit %d", h.entityType, results, name, limit))

	w.Header().Set("Content-Type", "application/json")
	err = h.encoder.Encode(w, results)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error encoding searched %s %v: %v", h.entityType, results, err.Error()))
		h.writeUnexpectedError(w)
		return
	}
}

func (h *SearchHandler[T]) readLimit(r *http.Request) (int, error) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return h.defaultLimit, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, fmt.Errorf("limit must be an integer, found %s", limitStr)
	}

	if limit < 1 || limit > h.maxLimit {
		return 0, fmt.Errorf("limit must be in interval [1, %d]", h.maxLimit)
	}

	return limit, nil
}

func (h *SearchHandler[T]) writeUnexpectedError(w http.ResponseWriter) {
	http.Error(
		w,
		fmt.Sprintf("Unexpected error: could not perform %s search", h.entityType),
		http.StatusInternalServerError,
	)
}

func (h *SearchHandler[T]) SetEncoder(encoder serialization.Encoder[[]T]) {
	h.encoder = encoder
}

func (h *SearchHandler[T]) SetMaxLimit(limit int) {
	h.maxLimit = limit
}

func (h *SearchHandler[T]) SetDefaultLimit(limit int) {
	h.defaultLimit = limit
}
