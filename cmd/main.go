package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	playlisthandler "festwrap/cmd/handler/playlist"
	"festwrap/cmd/handler/search"
	"festwrap/cmd/middleware"
	spotifyArtists "festwrap/internal/artist/spotify"
	"festwrap/internal/env"
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

func GetEnvWithDefaultOrFail[T env.EnvValue](key string, defaultValue T) T {
	variable, err := env.GetEnvWithDefault[T](key, defaultValue)
	if err != nil {
		log.Fatalf("Could not read variable %s", key)
	}
	return variable
}

func GetEnvStringOrFail(key string) string {
	variable := os.Getenv(key)
	if variable == "" {
		log.Fatalf("Could not read variable %s", key)
	}
	return variable
}

func main() {

	port := GetEnvWithDefaultOrFail[string]("FESTWRAP_PORT", "8080")
	maxConnsPerHost := GetEnvWithDefaultOrFail[int]("FESTWRAP_MAX_CONNS_PER_HOST", 10)
	timeoutSeconds := GetEnvWithDefaultOrFail[int]("FESTWRAP_TIMEOUT_SECONDS", 5)
	setlistfmApiKey := GetEnvStringOrFail("FESTWRAP_SETLISTFM_APIKEY")
	maxSetlistFMNumSearchPages := GetEnvWithDefaultOrFail[int]("FESTWRAP_SETLISTFM_NUM_SEARCH_PAGES", 3)
	maxUpdateArtists := GetEnvWithDefaultOrFail[int]("FESTWRAP_MAX_UPDATE_ARTISTS", 5)
	addSetlistSleepMs := GetEnvWithDefaultOrFail[int]("FESTWRAP_ADD_SETLIST_SLEEP_MS", 550)

	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger := logging.NewBaseLogger(slogLogger)

	httpClient := &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: maxConnsPerHost},
		Timeout:   time.Duration(timeoutSeconds) * time.Second,
	}
	baseHttpClient := httpclient.NewBaseHTTPClient(httpClient)
	httpSender := httpsender.NewBaseHTTPRequestSender(&baseHttpClient)

	mux := mux.NewRouter()
	mux.Use(middleware.NewAuthTokenExtractor().Middleware)

	artistRepository := spotifyArtists.NewSpotifyArtistRepository(&httpSender)
	artistSearcher := search.NewFunctionSearcher(artistRepository.SearchArtist)
	searchArtistsHandler := search.NewSearchHandler(&artistSearcher, "artists", logger)
	mux.HandleFunc("/artists/search", searchArtistsHandler.ServeHTTP).Methods(http.MethodGet)

	playlistRepository := spotifyplaylists.NewSpotifyPlaylistRepository(&httpSender)
	playlistSearcher := search.NewFunctionSearcher(playlistRepository.SearchPlaylist)
	userRepository := spotifyusers.NewSpotifyUserRepository(&httpSender)
	searchPlaylistsHandler := search.NewSearchHandler(&playlistSearcher, "playlists", logger)
	mux.HandleFunc(
		"/playlists/search",
		middleware.NewUserIdMiddleware(&searchPlaylistsHandler, userRepository).ServeHTTP,
	).Methods(http.MethodGet)

	setlistRepository := setlistfm.NewSetlistFMSetlistRepository(setlistfmApiKey, &httpSender)
	setlistRepository.SetMaxPages(maxSetlistFMNumSearchPages)
	songRepository := spotifysongs.NewSpotifySongRepository(&httpSender)
	playlistService := playlist.NewConcurrentPlaylistService(
		&playlistRepository,
		setlistRepository,
		songRepository,
	)
	existingPlaylistUpdateHandler := playlisthandler.NewUpdateExistingPlaylistHandler("playlistId", &playlistService, logger)
	existingPlaylistUpdateHandler.SetMaxArtists(maxUpdateArtists)
	existingPlaylistUpdateHandler.SetAddSetlistSleep(addSetlistSleepMs)
	mux.HandleFunc("/playlists/{playlistId}", existingPlaylistUpdateHandler.ServeHTTP).
		Name("playlistId").
		Methods(http.MethodPut)

	newPlaylistUpdateHandler := playlisthandler.NewUpdateNewPlaylistHandler(&playlistService, logger)
	newPlaylistUpdateHandler.SetMaxArtists(maxUpdateArtists)
	newPlaylistUpdateHandler.SetAddSetlistSleep(addSetlistSleepMs)
	mux.HandleFunc(
		"/playlists",
		middleware.NewUserIdMiddleware(&newPlaylistUpdateHandler, userRepository).ServeHTTP,
	).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	server.ListenAndServe()
}
