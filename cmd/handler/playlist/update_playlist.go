package playlist

import (
	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	builders "festwrap/internal/playlist/update_builders"
	"fmt"
	"net/http"
)

type UpdatePlaylistHandler struct {
	playlistService       playlist.PlaylistService
	logger                logging.Logger
	playlistUpdateBuilder playlist.PlaylistUpdateBuilder
	maxArtists            int
}

func NewUpdatePlaylistHandler(
	playlistService playlist.PlaylistService,
	playlistUpdateBuilder playlist.PlaylistUpdateBuilder,
	logger logging.Logger,
) UpdatePlaylistHandler {
	return UpdatePlaylistHandler{
		playlistService:       playlistService,
		logger:                logger,
		playlistUpdateBuilder: playlistUpdateBuilder,
		maxArtists:            5,
	}
}

func NewUpdateExistingPlaylistHandler(
	pathId string,
	playlistService playlist.PlaylistService,
	logger logging.Logger,
) UpdatePlaylistHandler {
	builder := builders.NewExistingPlaylistUpdateBuilder(pathId)
	return UpdatePlaylistHandler{
		playlistService:       playlistService,
		logger:                logger,
		playlistUpdateBuilder: &builder,
		maxArtists:            5,
	}
}

func NewUpdateNewPlaylistHandler(
	playlistService playlist.PlaylistService,
	logger logging.Logger,
) UpdatePlaylistHandler {
	builder := builders.NewNewPlaylistUpdateBuilder(playlistService)
	return UpdatePlaylistHandler{
		playlistService:       playlistService,
		logger:                logger,
		playlistUpdateBuilder: &builder,
		maxArtists:            5,
	}
}

func (h *UpdatePlaylistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	update, err := h.playlistUpdateBuilder.Build(r)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("could not get playlist update details: %v", err))
		http.Error(w, "could not obtain playlist details from request", http.StatusBadRequest)
		return
	}

	if len(update.Artists) == 0 || len(update.Artists) > h.maxArtists {
		message := fmt.Sprintf("validation error: number of artists must be between 1 and %d", h.maxArtists)
		h.logger.Warn(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	errors := 0
	for _, artist := range update.Artists {
		err := h.playlistService.AddSetlist(r.Context(), update.PlaylistId, artist.Name)
		if err != nil {
			message := fmt.Sprintf("could not add songs for %s to playlist %s: %v", artist.Name, update.PlaylistId, err)
			h.logger.Warn(message)
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

func (h *UpdatePlaylistHandler) SetPlaylistUpdateBuilder(builder playlist.PlaylistUpdateBuilder) {
	h.playlistUpdateBuilder = builder
}

func (h *UpdatePlaylistHandler) GetPlaylistService() playlist.PlaylistService {
	return h.playlistService
}

func (h *UpdatePlaylistHandler) SetPlaylistService(service playlist.PlaylistService) {
	h.playlistService = service
}

func (h *UpdatePlaylistHandler) SetMaxArtists(limit int) {
	h.maxArtists = limit
}
