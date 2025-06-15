package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	playlisthandler "festwrap/cmd/handler/playlist"
	"festwrap/cmd/handler/search"
	"festwrap/cmd/middleware"
	spotifyArtists "festwrap/internal/artist/spotify"
	httpclient "festwrap/internal/http/client"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	spotifyplaylists "festwrap/internal/playlist/spotify"
	"festwrap/internal/setlist/setlistfm"
	spotifysongs "festwrap/internal/song/spotify"
	spotifyusers "festwrap/internal/user/spotify"

	"github.com/gorilla/mux"
)

func setupLogger() logging.Logger {
	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logging.NewBaseLogger(slogLogger)
}

func setupHTTPSender(config Config) httpsender.HTTPRequestSender {
	httpClient := &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: config.MaxConnsPerHost},
		Timeout:   time.Duration(config.HttpClientTimeoutSeconds) * time.Second,
	}
	baseHttpClient := httpclient.NewBaseHTTPClient(httpClient)
	sender := httpsender.NewBaseHTTPRequestSender(&baseHttpClient)
	return &sender
}

func main() {

	config := ReadConfig()
	logger := setupLogger()
	httpSender := setupHTTPSender(config)

	mux := mux.NewRouter()
	mux.Use(middleware.NewAuthTokenExtractor().Middleware)

	// Set search artist endpoint
	artistRepository := spotifyArtists.NewSpotifyArtistRepository(httpSender)
	artistSearcher := search.NewFunctionSearcher(artistRepository.SearchArtist)
	searchArtistsHandler := search.NewSearchHandler(&artistSearcher, "artists", logger)
	mux.HandleFunc("/artists/search", searchArtistsHandler.ServeHTTP).Methods(http.MethodGet)

	// Set create new playlist endpoint
	playlistRepository := spotifyplaylists.NewSpotifyPlaylistRepository(httpSender)
	setlistRepository := setlistfm.NewSetlistFMSetlistRepository(config.SetlistfmApiKey, httpSender)
	setlistRepository.SetMaxPages(config.MaxSetlistFMNumSearchPages)
	songRepository := spotifysongs.NewSpotifySongRepository(httpSender)
	playlistService := playlist.NewConcurrentPlaylistService(
		&playlistRepository,
		setlistRepository,
		songRepository,
	)
	newPlaylistUpdateHandler := playlisthandler.NewUpdateNewPlaylistHandler(&playlistService, logger)
	newPlaylistUpdateHandler.SetMaxArtists(config.MaxUpdateArtists)
	newPlaylistUpdateHandler.SetAddSetlistSleep(config.AddSetlistSleepMs)
	userRepository := spotifyusers.NewSpotifyUserRepository(httpSender)
	userIdExtractor := middleware.NewUserIdExtractor(userRepository)
	mux.Handle(
		"/playlists",
		userIdExtractor.Middleware(http.HandlerFunc(newPlaylistUpdateHandler.ServeHTTP))).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: mux,
	}

	server.ListenAndServe()
}
