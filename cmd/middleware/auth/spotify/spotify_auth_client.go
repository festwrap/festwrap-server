package spotify

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
)

type SpotifyAuthClient struct {
	sender                  httpsender.HTTPRequestSender
	accessTokenDeserializer serialization.Deserializer[SpotifyAccessTokenInfo]
	clientId                string
	clientSecret            string
	host                    string
	refreshToken            string
	expirationTime          time.Time
	latestAccessToken       string
}

type SpotifyAccessTokenInfo struct {
	AccessToken           string `json:"access_token"`
	ExpirationTimeSeconds int    `json:"expires_in"`
}

func NewSpotifyAuthClient(
	sender httpsender.HTTPRequestSender,
	refreshToken string,
	clientId string,
	clientSecret string,
) SpotifyAuthClient {
	deserializer := serialization.NewJsonDeserializer[SpotifyAccessTokenInfo]()
	return SpotifyAuthClient{
		sender:                  sender,
		accessTokenDeserializer: &deserializer,
		clientId:                clientId,
		clientSecret:            clientSecret,
		refreshToken:            refreshToken,
		expirationTime:          time.Now(),
		host:                    "accounts.spotify.com/api/token",
	}
}

func (c *SpotifyAuthClient) GetAccessToken() (string, error) {
	now := time.Now()
	if now.Before(c.expirationTime) {
		return c.latestAccessToken, nil
	}

	tokenInfo, err := c.requestNewAccessToken()
	if err != nil {
		return "", fmt.Errorf("could not request new access token: %v", err)
	}

	c.latestAccessToken = tokenInfo.AccessToken
	newExpirationTime := time.Now().Add(time.Second * time.Duration(tokenInfo.ExpirationTimeSeconds))
	c.expirationTime = newExpirationTime
	return c.latestAccessToken, nil
}

func (c SpotifyAuthClient) requestNewAccessToken() (SpotifyAccessTokenInfo, error) {
	var accessTokenInfo SpotifyAccessTokenInfo
	accessTokenOpts, err := c.buildAccessTokenOpts()
	if err != nil {
		return accessTokenInfo, fmt.Errorf("could not build access token options: %v", err)
	}

	responseBody, err := c.sender.Send(accessTokenOpts)
	if err != nil {
		return accessTokenInfo, fmt.Errorf("error requesting access token: %v", err)
	}

	err = c.accessTokenDeserializer.Deserialize(*responseBody, &accessTokenInfo)
	if err != nil {
		return accessTokenInfo, fmt.Errorf("could not deserialize access token response")
	}

	return accessTokenInfo, nil
}

func (c SpotifyAuthClient) buildAccessTokenOpts() (httpsender.HTTPRequestOptions, error) {
	queryParams := url.Values{}
	queryParams.Set("grant_type", "refresh_token")
	queryParams.Set("refresh_token", c.refreshToken)
	queryParams.Set("client_id", c.clientId)
	url := fmt.Sprintf("https://%s?%s", c.host, queryParams.Encode())

	accessTokenOpts := httpsender.NewHTTPRequestOptions(url, httpsender.POST, http.StatusOK)

	authCode := base64.StdEncoding.EncodeToString([]byte(c.clientId + ":" + c.clientSecret))
	accessTokenOpts.SetHeaders(map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Basic " + authCode,
	})
	return accessTokenOpts, nil
}
