package spotify

import (
	"festwrap/internal/artist"
	"fmt"
)

type spotifyImage struct {
	Url string `json:"url"`
}

type spotifyArtist struct {
	Name   string         `json:"name"`
	Images []spotifyImage `json:"images"`
}

func (a spotifyArtist) GetSmallestImageUri() (string, error) {
	nImages := len(a.Images)
	if nImages == 0 {
		return "", fmt.Errorf("could not find image for artist %s", a.Name)
	}
	return a.Images[nImages-1].Url, nil
}

type spotifyArtists struct {
	ArtistItems []spotifyArtist `json:"items"`
}

type spotifyResponse struct {
	Artists spotifyArtists `json:"artists"`
}

func (s *spotifyResponse) GetArtists() []artist.Artist {
	result := []artist.Artist{}
	for _, currentArtist := range s.Artists.ArtistItems {
		imageUri, err := currentArtist.GetSmallestImageUri()
		resultArtist := artist.NewArtist(currentArtist.Name)
		if err == nil {
			resultArtist.SetImageUri(imageUri)
		}
		result = append(result, resultArtist)
	}
	return result
}
