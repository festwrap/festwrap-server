"use client"
import querystring from "querystring"
import Image from "next/image"
import {
  redirectToSpotifyAuth,
  requestSpotifyToken,
  SpotifyAuthConfig,
  SpotifyAuthTokens,
} from "./Authorization"
import { SpotifyCredentials } from "./Credentials"
import { REDIRECT_URI, SPOTIFY_CLIENT_ID, SPOTIFY_URL } from "../../env"

const AUTH_URL = `${SPOTIFY_URL}authorize?`
const SCOPE =
  "playlist-modify-private playlist-modify-public playlist-read-private"

export default function Home() {
  const authorizeSpotify = () => {
    if (!SPOTIFY_CLIENT_ID) {
      throw new Error("Spotify client ID not found")
    }

    const queryString = querystring.stringify({
      response_type: "code",
      client_id: SPOTIFY_CLIENT_ID,
      scope: SCOPE,
      redirect_uri: REDIRECT_URI,
    })

    const authUrl = AUTH_URL + queryString

    window.location.href = authUrl
  }

  const copyToClipboardToken = (token: SpotifyAuthTokens) => {
    navigator.clipboard.writeText(localStorage.getItem(token) || "")
  }

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="flex flex-col items-center space-y-8">
        <h1 className="text-4xl font-bold">Festwrap</h1>
        <button
          className="bg-green-500 text-white px-4 py-2 rounded-md"
          onClick={() => authorizeSpotify()}
        >
          Login with Spotify
        </button>
        <div className="flex flex-row gap-2">
          <button
            className="border border-gray-200 text-gray-800 px-4 py-2 rounded-md"
            onClick={() => copyToClipboardToken(SpotifyAuthTokens.ACCESS_TOKEN)}
          >
            Copy access token
          </button>
          <button
            className="border border-gray-200 text-gray-800 px-4 py-2 rounded-md"
            onClick={() =>
              copyToClipboardToken(SpotifyAuthTokens.REFRESH_TOKEN)
            }
          >
            Copy refresh token
          </button>
        </div>
      </div>
    </main>
  )
}
