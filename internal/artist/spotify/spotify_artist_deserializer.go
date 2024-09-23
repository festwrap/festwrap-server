package spotify

import (
	"encoding/json"

	"festwrap/internal/artist"
	"festwrap/internal/serialization/errors"
)

type SpotifyArtistDeserializer struct{}

func NewSpotifyArtistDeserializer() SpotifyArtistDeserializer {
	return SpotifyArtistDeserializer{}
}

func (s *SpotifyArtistDeserializer) Deserialize(setlist []byte) (*[]artist.Artist, error) {
	var response SpotifyResponse
	err := json.Unmarshal(setlist, &response)
	if err != nil {
		return nil, errors.NewDeserializationError(err.Error())
	}

	result := []artist.Artist{}
	for _, currentArtist := range response.Artists.ArtistItems {
		imageUri, err := currentArtist.GetSmallestImageUri()
		resultArtist := artist.NewArtist(currentArtist.Name)
		if err == nil {
			resultArtist.SetImageUri(imageUri)
		}
		result = append(result, resultArtist)
	}
	return &result, nil
}
