package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"festwrap/internal/artist"
	"festwrap/internal/serialization"
)

type SearchArtistHandler struct {
	encoder      serialization.Encoder[[]artist.Artist]
	defaultLimit int
	maxLimit     int
	repository   artist.ArtistRepository
}

func NewSearchArtistHandler(repository artist.ArtistRepository) SearchArtistHandler {
	return SearchArtistHandler{
		encoder:      serialization.NewJsonEncoder[[]artist.Artist](),
		repository:   repository,
		maxLimit:     10,
		defaultLimit: 5,
	}
}

func (h *SearchArtistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Validation error: artist name must be provided", http.StatusBadRequest)
		return
	}

	limit, err := h.readLimit(r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Validation error: %v", err.Error()),
			http.StatusUnprocessableEntity,
		)
		return
	}

	artists, err := h.repository.SearchArtist(r.Context(), name, limit)
	if err != nil {
		h.writeUnexpectedError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = h.encoder.Encode(w, artists)
	if err != nil {
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
