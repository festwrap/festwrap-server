package playlist

import (
	services "festwrap/cmd/services"
	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	"festwrap/internal/serialization"
	"fmt"
	"io"
	"net/http"
)

type CreatePlaylistHandler struct {
	playlistService     services.PlaylistService
	logger              logging.Logger
	maxArtists          int
	maxArtistNameLength int
	requestDeserializer serialization.Deserializer[NewPlaylistRequest]
	responseEncoder     serialization.Encoder[CreatePlaylistResponse]
}

func NewCreatePlaylistHandler(
	playlistService services.PlaylistService,
	logger logging.Logger,
) CreatePlaylistHandler {
	requestDeserializer := serialization.NewJsonDeserializer[NewPlaylistRequest]()
	responseEncoder := serialization.NewJsonEncoder[CreatePlaylistResponse]()
	return CreatePlaylistHandler{
		playlistService:     playlistService,
		logger:              logger,
		maxArtists:          5,
		maxArtistNameLength: 50,
		requestDeserializer: &requestDeserializer,
		responseEncoder:     &responseEncoder,
	}
}

func (h *CreatePlaylistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("could not read create playlist body from request: %v", err))
		http.Error(w, "could not read body from request", http.StatusBadRequest)
		return
	}

	var newPlaylistRequest NewPlaylistRequest
	err = h.requestDeserializer.Deserialize(requestBody, &newPlaylistRequest)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to deserialize playlist information: %v", err))
		http.Error(w, "failed to read playlist information", http.StatusBadRequest)
		return
	}

	artists := newPlaylistRequest.Artists
	h.logger.Info(fmt.Sprintf("creating playlist with artists: %v", artists))
	if len(artists) == 0 || len(artists) > h.maxArtists {
		message := fmt.Sprintf("validation error: number of artists must be between 1 and %d", h.maxArtists)
		h.logger.Warn(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	for _, artist := range artists {
		if len(artist.Name) > h.maxArtistNameLength || len(artist.Name) == 0 {
			message := fmt.Sprintf(
				"validation error: artist name '%s' length should be in interval [1, %d]",
				artist,
				h.maxArtistNameLength,
			)
			h.logger.Warn(message)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
	}

	artistNames := newPlaylistRequest.GetArtistNames()
	result, err := h.playlistService.CreatePlaylistWithArtists(
		r.Context(),
		playlist.PlaylistDetails{Name: newPlaylistRequest.Playlist.Name, Description: "", IsPublic: true},
		artistNames,
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("could not create playlist :%v", err))
		http.Error(w, "unexpected error, could not create playlist", http.StatusInternalServerError)
		return
	}

	statusCode := http.StatusCreated
	switch result.Status {
	case services.Success:
		statusCode = http.StatusCreated
	case services.PartialFailure:
		statusCode = http.StatusMultiStatus
	default:
		h.logger.Warn(fmt.Sprintf("unexpected creation status %v", result.Status))
		http.Error(w, "unexpected error, could not create playlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)

	message := fmt.Sprintf("created playlist with id %s and artists %v", result.PlaylistId, artists)
	h.logger.Info(message)

	response := CreatePlaylistResponse{Playlist: CreatedPlaylist{Id: result.PlaylistId}}
	if err = h.responseEncoder.Encode(w, response); err != nil {
		message := fmt.Sprintf("encoding error, could not encode response: %v", err)
		h.logger.Error(message)
		http.Error(w, "unexpected error, could not encode response", http.StatusInternalServerError)
		return
	}
}

func (h *CreatePlaylistHandler) GetPlaylistService() services.PlaylistService {
	return h.playlistService
}

func (h *CreatePlaylistHandler) SetPlaylistService(service services.PlaylistService) {
	h.playlistService = service
}

func (h *CreatePlaylistHandler) SetMaxArtists(limit int) {
	h.maxArtists = limit
}

func (h *CreatePlaylistHandler) SetMaxArtistNameLength(length int) {
	h.maxArtistNameLength = length
}
