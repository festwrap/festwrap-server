package playlist

import (
	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	builders "festwrap/internal/playlist/update_builders"
	"festwrap/internal/serialization"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TODO: to separate file?
type PlaylistArtist struct {
	Name string `json:"name"`
}

type NewPlaylist struct {
	Name string `json:"name"`
}

type NewPlaylistRequest struct {
	Playlist playlist.Playlist `json:"playlist"`
	Artists  []PlaylistArtist  `json:"artists"`
}

type CreatedPlaylist struct {
	Id string `json:"id"`
}

type CreatePlaylistResponse struct {
	Playlist CreatedPlaylist `json:"playlist"`
}

type UpdatePlaylistHandler struct {
	playlistService       playlist.PlaylistService
	logger                logging.Logger
	playlistUpdateBuilder playlist.PlaylistUpdateBuilder
	maxArtists            int
	returnResponse        bool
	requestDeserializer   serialization.Deserializer[NewPlaylistRequest]
	responseEncoder       serialization.Encoder[CreatePlaylistResponse]
	successStatusCode     int
	addSetlistSleepMs     int
}

func NewUpdatePlaylistHandler(
	playlistService playlist.PlaylistService,
	playlistUpdateBuilder playlist.PlaylistUpdateBuilder,
	successStatusCode int,
	logger logging.Logger,
) UpdatePlaylistHandler {
	deserializer := serialization.NewJsonDeserializer[NewPlaylistRequest]()
	responseEncoder := serialization.NewJsonEncoder[CreatePlaylistResponse]()
	return UpdatePlaylistHandler{
		playlistService:       playlistService,
		logger:                logger,
		playlistUpdateBuilder: playlistUpdateBuilder,
		maxArtists:            5,
		returnResponse:        false,
		requestDeserializer:   &deserializer,
		responseEncoder:       &responseEncoder,
		successStatusCode:     successStatusCode,
		addSetlistSleepMs:     0,
	}
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
	r.Body.Close()
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("could read body from request: %v", err))
		http.Error(w, "could read body from request", http.StatusBadRequest)
		return
	}

	var newPlaylistRequest NewPlaylistRequest
	err = h.requestDeserializer.Deserialize(requestBody, &newPlaylistRequest)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to deserialize playlist information: %v", err))
		http.Error(w, "failed to deserialize playlist information", http.StatusBadRequest)
		return
	}

	artists := newPlaylistRequest.Artists
	if len(artists) == 0 || len(artists) > h.maxArtists {
		message := fmt.Sprintf("validation error: number of artists must be between 1 and %d", h.maxArtists)
		h.logger.Warn(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	playlist := playlist.PlaylistDetails{

	}
	h.playlistService.CreatePlaylistWithArtists(r.Context(), )
	errors := 0
	for i, artist := range update.Artists {
		if i > 0 {
			// Sleep to avoid hitting Setlistfm rate limit
			time.Sleep(time.Duration(h.addSetlistSleepMs) * time.Millisecond)
		}

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

func (h *UpdatePlaylistHandler) SetAddSetlistSleep(sleepMs int) {
	h.addSetlistSleepMs = sleepMs
}
