import querystring from "querystring"
import axios, { AxiosResponse } from "axios"
import crypto from "crypto"

export class SpotifyAuthConfig {
  constructor(
    public authUrl: string,
    public tokenUrl: string,
    public redirectUri: string,
    public scope: string
  ) {}

  getAuthUrl(): string {
    return this.authUrl
  }

  getTokenUrl(): string {
    return this.tokenUrl
  }

  getRedirectUri(): string {
    return this.redirectUri
  }

  getScope(): string {
    return this.scope
  }
}

const generateRandomString = (length: number): string => {
  const possible =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
  const values = crypto.randomBytes(length)
  return Array.from(values)
    .map((byte) => possible[byte % possible.length])
    .join("")
}

export function redirectToSpotifyAuth(
  response: any,
  authConfig: SpotifyAuthConfig,
  clientId: string,
  randomStringLength: number
): void {
  const state = generateRandomString(randomStringLength)
  response.redirect(
    authConfig.getAuthUrl() +
      querystring.stringify({
        response_type: "code",
        client_id: clientId,
        scope: authConfig.getScope(),
        redirect_uri: authConfig.getRedirectUri(),
        state: state,
      })
  )
}

export async function requestSpotifyToken(
  authConfig: SpotifyAuthConfig,
  authCode: string,
  secretBase64: string
): Promise<string> {
  const tokenRequestOptions = {
    method: "post",
    url: authConfig.getTokenUrl(),
    data: querystring.stringify({
      code: authCode,
      grant_type: "client_credentials",
    }),
    headers: {
      "content-type": "application/x-www-form-urlencoded",
      Authorization: "Basic " + secretBase64,
    },
  }

  return axios(tokenRequestOptions)
    .then((authResponse: AxiosResponse) => {
      return authResponse.data.access_token
    })
    .catch((error: any) => {
      throw new Error(`Could not get access token: ${error}`)
    })
}
