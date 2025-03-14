package spotify

import (
	"context"
	"errors"
	"fmt"

	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
)

type SpotifyUserRepository struct {
	tokenKey     types.ContextKey
	deserializer serialization.Deserializer[spotifyUserResponse]
	httpSender   httpsender.HTTPRequestSender
	host         string
}

func NewSpotifyUserRepository(httpSender httpsender.HTTPRequestSender) SpotifyUserRepository {
	return SpotifyUserRepository{
		tokenKey:     "token",
		deserializer: serialization.NewJsonDeserializer[spotifyUserResponse](),
		httpSender:   httpSender,
		host:         "api.spotify.com",
	}
}

func (r SpotifyUserRepository) GetCurrentUserId(ctx context.Context) (string, error) {
	token, ok := ctx.Value(r.tokenKey).(string)
	if !ok {
		return "", errors.New("could not retrieve token from context")
	}

	responseBody, err := r.httpSender.Send(r.getCurrentUserIdHTTPOptions(token))
	if err != nil {
		return "", fmt.Errorf("could not get current user: %v", err.Error())
	}

	var response spotifyUserResponse
	err = r.deserializer.Deserialize(*responseBody, &response)
	if err != nil {
		return "", fmt.Errorf("deserialization error: %v", err.Error())
	}

	return response.UserId, nil
}

func (r SpotifyUserRepository) getCurrentUserIdHTTPOptions(accessToken string) httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://%s/v1/me", r.host)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
	)
	return httpOptions
}

func (r *SpotifyUserRepository) SetTokenKey(key types.ContextKey) {
	r.tokenKey = key
}
