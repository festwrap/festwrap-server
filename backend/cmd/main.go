package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"festwrap/cmd/handler"
	"festwrap/cmd/middleware"
	"festwrap/internal/artist/spotify"
	"festwrap/internal/env"
	httpclient "festwrap/internal/http/client"
	httpsender "festwrap/internal/http/sender"
)

func GetEnvWithDefaultOrFail[T env.EnvValue](key string, defaultValue T) T {
	variable, err := env.GetEnvWithDefault[T](key, defaultValue)
	if err != nil {
		log.Fatalf("Could not read variable %s", key)
	}
	return variable
}

func main() {

	port := GetEnvWithDefaultOrFail[string]("FESTWRAP_PORT", "8080")
	maxConnsPerHost := GetEnvWithDefaultOrFail[int]("FESTWRAP_MAX_CONNS_PER_HOST", 10)
	timeoutSeconds := GetEnvWithDefaultOrFail[int]("FESTWRAP_TIMEOUT_SECONDS", 5)

	httpClient := &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: maxConnsPerHost},
		Timeout:   time.Duration(timeoutSeconds) * time.Second,
	}
	baseHttpClient := httpclient.NewBaseHTTPClient(httpClient)
	httpSender := httpsender.NewBaseHTTPRequestSender(&baseHttpClient)

	mux := http.NewServeMux()

	repository := spotify.NewSpotifyArtistRepository(&httpSender)
	searchArtistsHandler := handler.NewSearchArtistHandler(repository)
	mux.HandleFunc("/artists/search", searchArtistsHandler.ServeHTTP)

	wrappedMux := middleware.NewAuthTokenMiddleware(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: wrappedMux,
	}

	server.ListenAndServe()
}
