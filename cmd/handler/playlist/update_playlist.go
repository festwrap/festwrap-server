package playlist

import (
	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	"festwrap/internal/serialization"
	"fmt"
	"io"
	"net/http"
)

// TODO fix tests

// TODO move to separate file
type CreatePlaylistArtist struct {
	Name string `json:"name"`
}

type NewPlaylist struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

func (p NewPlaylist) toPlaylist() playlist.Playlist {
	return playlist.Playlist{
		Name:        p.Name,
		Description: p.Description,
		IsPublic:    p.IsPublic,
	}
}

type NewPlaylistRequest struct {
	Playlist NewPlaylist            `json:"playlist"`
	Artists  []CreatePlaylistArtist `json:"artists"`
}

type CreatedPlaylist struct {
	Id string `json:"id"`
}

type CreatePlaylistResponse struct {
	Playlist CreatedPlaylist `json:"playlist"`
}

type UpdatePlaylistHandler struct {
	playlistService     playlist.PlaylistService
	logger              logging.Logger
	maxArtists          int
	requestDeserializer serialization.Deserializer[NewPlaylistRequest]
	responseEncoder     serialization.Encoder[CreatePlaylistResponse]
	addSetlistSleepMs   int
}

func NewUpdatePlaylistHandler(
	playlistService playlist.PlaylistService,
	playlistUpdateBuilder playlist.PlaylistUpdateBuilder,
	logger logging.Logger,
) UpdatePlaylistHandler {
	deserializer := serialization.NewJsonDeserializer[NewPlaylistRequest]()
	responseEncoder := serialization.NewJsonEncoder[CreatePlaylistResponse]()
	return UpdatePlaylistHandler{
		playlistService:     playlistService,
		logger:              logger,
		maxArtists:          5,
		requestDeserializer: &deserializer,
		responseEncoder:     &responseEncoder,
		addSetlistSleepMs:   0,
	}
}

func (h *UpdatePlaylistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO simplify handler
	defer r.Body.Close()
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("could not read body from request: %v", err))
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

	if len(newPlaylistRequest.Artists) == 0 || len(newPlaylistRequest.Artists) > h.maxArtists {
		message := fmt.Sprintf("validation error: number of artists must be between 1 and %d", h.maxArtists)
		h.logger.Warn(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	// TODO move to method
	artists := make([]playlist.PlaylistArtist, len(newPlaylistRequest.Artists))
	for i, artist := range newPlaylistRequest.Artists {
		artists[i] = playlist.PlaylistArtist{Name: artist.Name}
	}
	result, err := h.playlistService.CreatePlaylistWithArtists(
		r.Context(), newPlaylistRequest.Playlist.toPlaylist(), artists,
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("could not create playlist: %v", err))
		http.Error(w, "could not add any artist to playlist", http.StatusInternalServerError)
		return
	}

	// TODO move to separate method
	var statusCode int
	if result.Status == playlist.PartialFailure {
		statusCode = http.StatusMultiStatus
		message := fmt.Sprintf(
			"partial failure: could not add some artists in %v to playlist: %v", newPlaylistRequest.Artists, err,
		)
		h.logger.Warn(message)
	} else if result.Status == playlist.Success {
		statusCode = http.StatusCreated
		h.logger.Info(
			fmt.Sprintf("added artists %v to playlist with id: %s", newPlaylistRequest.Artists, result.PlaylistId),
		)
	} else {
		h.logger.Error(fmt.Sprintf("unexpected status value: %v", result.Status))
		http.Error(w, "unexpected error: could not add any artist to playlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)

	response := CreatePlaylistResponse{Playlist: CreatedPlaylist{Id: result.PlaylistId}}
	if err = h.responseEncoder.Encode(w, response); err != nil {
		message := fmt.Sprintf("encoding error: could not encode response: %v", err)
		h.logger.Error(message)
		http.Error(w, "unexpected error: could not encode response", http.StatusInternalServerError)
		return
	}
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

func (h *UpdatePlaylistHandler) SetAddSetlistSleep(sleepMs int) {
	h.addSetlistSleepMs = sleepMs
}
