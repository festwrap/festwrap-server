package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"festwrap/internal/playlist"
	spotifyPlaylist "festwrap/internal/playlist/spotify"
	"festwrap/internal/setlist/setlistfm"
	spotifySong "festwrap/internal/song/spotify"
)

func main() {
	spotifyAccessToken := flag.String("spotify-token", "", "Spotify access token")
	setlistfmApiKey := flag.String("setlistfm-key", "", "Setlistfm API Key")
	artist := flag.String("artist", "", "Artist to add to the playlist")
	playlistId := flag.String("playlist-id", "", "Spotify playlist identifier where to add songs")
	spotifyHost := flag.String("spotify-api", "api.spotify.com", "Spotify API url")
	setlistfmHost := flag.String("setlistfm-api", "api.setlist.fm", "Setlistfm API url")
	minSongsPerSetlist := flag.Int("min-setlist-songs", 5, "Minimum number of songs to include")
	flag.Parse()

	httpClient := &http.Client{}

	fmt.Printf("Adding latest setlist songs for %s into Spotify playlist with id %s \n", *artist, *playlistId)

	setlistFmParser := setlistfm.NewSetlistFMParser()
	setlistFmParser.SetMinimumSongs(*minSongsPerSetlist)
	setlistRepository := setlistfm.NewSetlistFMSetlistRepository(
		&setlistfm.SetlistFMSetlistRepositoryConfig{Client: httpClient, Host: *setlistfmHost, ApiKey: *setlistfmApiKey},
		&setlistFmParser,
	)

	spotifySongParser := spotifySong.NewSpotifySongsParser()
	songRepository := spotifySong.NewSpotifySongRepository(
		httpClient,
		spotifySong.SpotifySongRepositoryConfig{Host: *spotifyHost, AccessToken: *spotifyAccessToken},
		&spotifySongParser,
	)

	playlistRepository := spotifyPlaylist.NewSpotifyPlaylistRepository(*spotifyHost, httpClient, *spotifyAccessToken)
	playlistService := playlist.NewConcurrentPlaylistService(
		&playlistRepository,
		setlistRepository,
		songRepository,
	)

	err := playlistService.AddSetlist(*playlistId, *artist)
	if err != nil {
		message := fmt.Sprintf("Could not add songs to setlist: %v", err)
		fmt.Println(message)
		os.Exit(1)
	}
}
