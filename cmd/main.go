package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	playlisthandler "festwrap/cmd/handler/playlist"
	"festwrap/cmd/handler/search"
	"festwrap/cmd/middleware"
	auth "festwrap/cmd/middleware/auth"
	spotifyauth "festwrap/cmd/middleware/auth/spotify"
	services "festwrap/cmd/services"
	spotifyArtists "festwrap/internal/artist/spotify"
	"festwrap/internal/event"
	httpclient "festwrap/internal/http/client"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/logging"
	"festwrap/internal/messaging"
	spotifyplaylists "festwrap/internal/playlist/spotify"
	"festwrap/internal/setlist/setlistfm"
	spotifysongs "festwrap/internal/song/spotify"
	spotifyusers "festwrap/internal/user/spotify"

	"cloud.google.com/go/pubsub"
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
	spotifyAuthClient := spotifyauth.NewSpotifyAuthClient(
		httpSender, config.SpotifyRefreshToken, config.SpotifyClientId, config.SpotifyClientSecret,
	)
	mux.Use(auth.NewAuthTokenExtractor(&spotifyAuthClient, logger).Middleware)

	// Configure pubsub client
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, config.PubsubProjectId)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initialize pubsub client: %s", err))
		os.Exit(1)
	}
	defer pubsubClient.Close()
	publisher := messaging.NewPubsubPublisher(pubsubClient, logger)

	// Set search artist endpoint
	artistRepository := spotifyArtists.NewSpotifyArtistRepository(httpSender)
	artistSearcher := search.NewFunctionSearcher(artistRepository.SearchArtist)
	searchArtistsHandler := search.NewSearchHandler(&artistSearcher, "artists", logger)
	searchArtistsHandler.SetMaxNameLength(config.MaxArtistNameLength)
	mux.HandleFunc("/artists/search", searchArtistsHandler.ServeHTTP).Methods(http.MethodGet)

	// Initialize playlist service
	playlistRepository := spotifyplaylists.NewSpotifyPlaylistRepository(httpSender)
	setlistRepository := setlistfm.NewSetlistFMSetlistRepository(config.SetlistfmApiKey, httpSender)
	setlistRepository.SetMaxPages(config.MaxSetlistFMNumSearchPages)
	setlistRepository.SetNextPageSleep(config.NextPageSleepMs)
	songRepository := spotifysongs.NewSpotifySongRepository(httpSender)
	playlistService := services.NewBasePlaylistService(
		&playlistRepository,
		setlistRepository,
		songRepository,
		logger,
	)
	playlistService.SetAddSetlistSleep(config.AddSetlistSleepMs)

	// Configure service to publish creation events
	createNotifier := event.NewBaseNotifier[event.PlaylistCreatedEvent]()
	publishObserver := event.NewPublishEventObserver[event.PlaylistCreatedEvent](publisher, config.CreatePlaylistTopic)
	createNotifier.AddObserver(publishObserver)
	playlistService.SetPlaylistCreateNotifier(createNotifier)

	// Set create new playlist endpoint
	newPlaylistUpdateHandler := playlisthandler.NewCreatePlaylistHandler(&playlistService, logger)
	newPlaylistUpdateHandler.SetMaxArtists(config.MaxCreateArtists)
	newPlaylistUpdateHandler.SetMaxArtistNameLength(config.MaxArtistNameLength)
	userRepository := spotifyusers.NewSpotifyUserRepository(httpSender)
	userIdExtractor := middleware.NewUserIdExtractor(userRepository, logger)
	mux.Handle(
		"/playlists",
		userIdExtractor.Middleware(http.HandlerFunc(newPlaylistUpdateHandler.ServeHTTP))).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: mux,
	}

	logger.Info(fmt.Sprintf("Starting server at port %s", config.Port))
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(fmt.Sprintf("could not start server %v", err))
	}
}
