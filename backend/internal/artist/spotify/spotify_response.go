package spotify

import (
	"festwrap/internal/artist/errors"
	"fmt"
)

type SpotifyImage struct {
	Url string `json:"url"`
}

type SpotifyArtist struct {
	Name   string         `json:"name"`
	Images []SpotifyImage `json:"images"`
}

func (a SpotifyArtist) GetSmallestImageUri() (string, error) {
	nImages := len(a.Images)
	if nImages == 0 {
		return "", errors.NewImageNotFoundError(fmt.Sprintf("Could not find image for artist %s", a.Name))
	}
	return a.Images[nImages-1].Url, nil
}

type SpotifyArtists struct {
	ArtistItems []SpotifyArtist `json:"items"`
}

type SpotifyResponse struct {
	Artists SpotifyArtists `json:"artists"`
}
