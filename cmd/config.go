package main

import (
	"festwrap/internal/env"
	"log"
	"os"
)

type Config struct {
	Port                       string
	MaxConnsPerHost            int
	TimeoutSeconds             int
	SetlistfmApiKey            string
	MaxSetlistFMNumSearchPages int
	MaxUpdateArtists           int
	AddSetlistSleepMs          int
}

func ReadConfig() Config {
	return Config{
		Port:                       GetEnvWithDefaultOrFail[string]("FESTWRAP_PORT", "8080"),
		MaxConnsPerHost:            GetEnvWithDefaultOrFail[int]("FESTWRAP_MAX_CONNS_PER_HOST", 10),
		TimeoutSeconds:             GetEnvWithDefaultOrFail[int]("FESTWRAP_TIMEOUT_SECONDS", 5),
		SetlistfmApiKey:            GetEnvStringOrFail("FESTWRAP_SETLISTFM_APIKEY"),
		MaxSetlistFMNumSearchPages: GetEnvWithDefaultOrFail[int]("FESTWRAP_SETLISTFM_NUM_SEARCH_PAGES", 3),
		MaxUpdateArtists:           GetEnvWithDefaultOrFail[int]("FESTWRAP_MAX_UPDATE_ARTISTS", 5),
		AddSetlistSleepMs:          GetEnvWithDefaultOrFail[int]("FESTWRAP_ADD_SETLIST_SLEEP_MS", 550),
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
