package playlist

import (
	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	"festwrap/internal/serialization"
	"fmt"
	"io"
	"net/http"
)

type UpdatePlaylistHandler struct {
	playlistService playlist.PlaylistService
	logger          logging.Logger
	deserializer    serialization.Deserializer[playlistUpdate]
	maxArtists      int
	playlistIdPath  string
}

func NewUpdatePlaylistHandler(playlistService playlist.PlaylistService, logger logging.Logger) UpdatePlaylistHandler {
	return UpdatePlaylistHandler{
		playlistService: playlistService,
		logger:          logger,
		deserializer:    serialization.NewJsonDeserializer[playlistUpdate](),
		maxArtists:      5,
		playlistIdPath:  "playlistId",
	}
}

func (h *UpdatePlaylistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	playlistId := r.PathValue(h.playlistIdPath)
	if playlistId == "" {
		message := "validation error: playlist id was not provided"
		h.logger.Warn(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		message := "validation error: could not read body"
		h.logger.Warn(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	var update playlistUpdate
	h.deserializer.Deserialize(requestBody, &update)
	if len(update.Artists) > h.maxArtists {
		message := fmt.Sprintf("validation error: cannot update playlist with more than %d artists", h.maxArtists)
		h.logger.Warn(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	errors := 0
	for _, artist := range update.Artists {
		err := h.playlistService.AddSetlist(r.Context(), playlistId, artist.Name)
		if err != nil {
			h.logger.Warn(fmt.Sprintf("could not add songs for %s to playlist %s: %v", artist.Name, playlistId, err))
			errors += 1
		}
	}

	statusCode := http.StatusCreated
	if errors > 0 && errors < len(update.Artists) {
		statusCode = http.StatusMultiStatus
	} else if errors > 0 {
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
}

func (h *UpdatePlaylistHandler) SetDeserializer(deserializer serialization.Deserializer[playlistUpdate]) {
	h.deserializer = deserializer
}

func (h *UpdatePlaylistHandler) SetPlaylistService(service playlist.PlaylistService) {
	h.playlistService = service
}

func (h *UpdatePlaylistHandler) SetMaxArtists(limit int) {
	h.maxArtists = limit
}
