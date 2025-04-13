package playlist

import (
	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	builders "festwrap/internal/playlist/update_builders"
	"festwrap/internal/serialization"
	"fmt"
	"net/http"
)

type Playlist struct {
	Id string `json:"id"`
}

type UpdatePlaylistResponse struct {
	Playlist Playlist `json:"playlist"`
}

type UpdatePlaylistHandler struct {
	playlistService       playlist.PlaylistService
	logger                logging.Logger
	playlistUpdateBuilder playlist.PlaylistUpdateBuilder
	maxArtists            int
	returnResponse        bool
	responseEncoder       serialization.Encoder[UpdatePlaylistResponse]
	successStatusCode     int
}

func NewUpdatePlaylistHandler(
	playlistService playlist.PlaylistService,
	playlistUpdateBuilder playlist.PlaylistUpdateBuilder,
	successStatusCode int,
	logger logging.Logger,
) UpdatePlaylistHandler {
	responseEncoder := serialization.NewJsonEncoder[UpdatePlaylistResponse]()
	return UpdatePlaylistHandler{
		playlistService:       playlistService,
		logger:                logger,
		playlistUpdateBuilder: playlistUpdateBuilder,
		maxArtists:            5,
		returnResponse:        false,
		responseEncoder:       &responseEncoder,
		successStatusCode:     successStatusCode,
	}
}

func NewUpdateExistingPlaylistHandler(
	pathId string,
	playlistService playlist.PlaylistService,
	logger logging.Logger,
) UpdatePlaylistHandler {
	builder := builders.NewExistingPlaylistUpdateBuilder(pathId)
	handler := NewUpdatePlaylistHandler(playlistService, &builder, http.StatusOK, logger)
	return handler
}

func NewUpdateNewPlaylistHandler(
	playlistService playlist.PlaylistService,
	logger logging.Logger,
) UpdatePlaylistHandler {
	builder := builders.NewNewPlaylistUpdateBuilder(playlistService)
	handler := NewUpdatePlaylistHandler(playlistService, &builder, http.StatusCreated, logger)
	handler.ReturnResponse(true)
	return handler
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

	statusCode := h.successStatusCode
	if errors > 0 && errors < len(update.Artists) {
		statusCode = http.StatusMultiStatus
	} else if errors > 0 {
		statusCode = http.StatusInternalServerError
	}
	w.WriteHeader(statusCode)

	message := fmt.Sprintf("updated playlist with id: %s and artists: %v", update.PlaylistId, update.Artists)
	h.logger.Info(message)

	if h.returnResponse {
		response := UpdatePlaylistResponse{Playlist: Playlist{Id: update.PlaylistId}}
		if err = h.responseEncoder.Encode(w, response); err != nil {
			message := fmt.Sprintf("encoding error: could not encode response: %v", err)
			h.logger.Error(message)
			http.Error(w, "unexpected error: could not encode response", http.StatusInternalServerError)
			return
		}
	}
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

func (h *UpdatePlaylistHandler) ReturnResponse(flag bool) {
	h.returnResponse = flag
}

func (h *UpdatePlaylistHandler) SetSuccessStatusCode(status int) {
	h.successStatusCode = status
}
