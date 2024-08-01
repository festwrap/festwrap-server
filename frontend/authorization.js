import querystring from 'querystring';
import axios from 'axios';
import crypto from 'crypto';

export class SpotifyAuthConfig {
  constructor(authUrl, tokenUrl, redirectUri, scope) {
    this.authUrl = authUrl;
    this.tokenUrl = tokenUrl;
    this.redirectUri = redirectUri;
    this.scope = scope;
  }

  getAuthUrl() {
    return this.authUrl;
  }

  getTokenUrl() {
    return this.tokenUrl
  }

  getRedirectUri() {
    return this.redirectUri;
  }

  getScope() {
    return this.scope;
  }

}

const generateRandomString = (length) => {
  const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const values = crypto.getRandomValues(new Uint8Array(length));
  return values.reduce((acc, x) => acc + possible[x % possible.length], "");
}

export function redirectToSpotifyAuth(response, authConfig, clientId, randomStringLength) {
  const state = generateRandomString(randomStringLength);
  response.redirect(authConfig.getAuthUrl() +
    querystring.stringify({
      response_type: 'code',
      client_id: clientId,
      scope: authConfig.getScope(),
      redirect_uri: authConfig.getRedirectUri(),
      state: state
    }));
}

export async function requestSpotifyToken(authConfig, authCode, secretBase64) {
  const tokenRequestOptions = {
    method: 'post',
    url: authConfig.getTokenUrl(),
    data: {
      code: authCode,
      redirect_uri: authConfig.getRedirectUri(),
      grant_type: 'authorization_code'
    },
    headers: {
      'content-type': 'application/x-www-form-urlencoded',
      'Authorization': 'Basic ' + secretBase64
    },
    json: true
  };

  return axios(tokenRequestOptions)
  .then((authResponse) => {
    return authResponse.data.access_token;
  })
  .catch((error) => {
    throw new Error(`Could not get access token: ${error}`);
  })
}
