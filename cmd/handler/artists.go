package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"festwrap/internal/artist"
	"festwrap/internal/logging"
	"festwrap/internal/serialization"
)

type SearchArtistHandler struct {
	encoder      serialization.Encoder[[]artist.Artist]
	defaultLimit int
	maxLimit     int
	repository   artist.ArtistRepository
	logger       logging.Logger
}

func NewSearchArtistHandler(repository artist.ArtistRepository, logger logging.Logger) SearchArtistHandler {
	return SearchArtistHandler{
		encoder:      serialization.NewJsonEncoder[[]artist.Artist](),
		repository:   repository,
		maxLimit:     10,
		defaultLimit: 5,
		logger:       logger,
	}
}

func (h *SearchArtistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		message := "Validation error: artist name was not provided"
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

	artists, err := h.repository.SearchArtist(r.Context(), name, limit)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error searching for artist: %v", err.Error()))
		h.writeUnexpectedError(w)
		return
	}
	h.logger.Info(fmt.Sprintf("Found artist %v for %s, using limit %d", artists, name, limit))

	w.Header().Set("Content-Type", "application/json")
	err = h.encoder.Encode(w, artists)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error encoding searched artists %v: %v", artists, err.Error()))
		h.writeUnexpectedError(w)
		return
	}
}

func (h *SearchArtistHandler) readLimit(r *http.Request) (int, error) {
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

func (h *SearchArtistHandler) writeUnexpectedError(w http.ResponseWriter) {
	http.Error(w, "Unexpected error: could not perform artist search", http.StatusInternalServerError)
}

func (h *SearchArtistHandler) SetEncoder(encoder serialization.Encoder[[]artist.Artist]) {
	h.encoder = encoder
}

func (h *SearchArtistHandler) SetMaxLimit(limit int) {
	h.maxLimit = limit
}

func (h *SearchArtistHandler) SetDefaultLimit(limit int) {
	h.defaultLimit = limit
}
