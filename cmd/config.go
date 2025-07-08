package main

import (
	"festwrap/internal/env"
	"log"
	"os"
)

type Config struct {
	Port                       string
	MaxConnsPerHost            int
	MaxSetlistFMNumSearchPages int
	MaxCreateArtists           int
	MaxArtistNameLength        int
	AddSetlistSleepMs          int
	NextPageSleepMs            int
	HttpClientTimeoutSeconds   int

	SetlistfmApiKey string

	SpotifyClientId     string
	SpotifyClientSecret string
	SpotifyRefreshToken string
}

func ReadConfig() Config {
	return Config{
		Port:                       GetEnvWithDefaultOrFail[string]("FESTWRAP_PORT", "8080"),
		MaxConnsPerHost:            GetEnvWithDefaultOrFail[int]("FESTWRAP_MAX_CONNS_PER_HOST", 10),
		SetlistfmApiKey:            GetEnvStringOrFail("FESTWRAP_SETLISTFM_APIKEY"),
		MaxSetlistFMNumSearchPages: GetEnvWithDefaultOrFail[int]("FESTWRAP_SETLISTFM_NUM_SEARCH_PAGES", 3),
		MaxCreateArtists:           GetEnvWithDefaultOrFail[int]("FESTWRAP_MAX_CREATE_ARTISTS", 5),
		MaxArtistNameLength:        GetEnvWithDefaultOrFail[int]("FESTWRAP_MAX_ARTIST_NAME_LENGTH", 50),
		AddSetlistSleepMs:          GetEnvWithDefaultOrFail[int]("FESTWRAP_ADD_SETLIST_SLEEP_MS", 550),
		NextPageSleepMs:            GetEnvWithDefaultOrFail[int]("FESTWRAP_GET_SETLIST_NEXT_PAGE_SLEEP_MS", 550),
		HttpClientTimeoutSeconds:   GetEnvWithDefaultOrFail[int]("FESTWRAP_HTTP_CLIENT_TIMEOUT_S", 5),
		SpotifyClientId:            GetEnvStringOrFail("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret:        GetEnvStringOrFail("SPOTIFY_CLIENT_SECRET"),
		SpotifyRefreshToken:        GetEnvStringOrFail("SPOTIFY_REFRESH_TOKEN"),
	}
}

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
