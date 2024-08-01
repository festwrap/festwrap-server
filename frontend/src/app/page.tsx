"use client"
import querystring from "querystring"
import Image from "next/image"
import {
  redirectToSpotifyAuth,
  requestSpotifyToken,
  SpotifyAuthConfig,
} from "./Authorization"
import { SpotifyCredentials } from "./Credentials"
import { REDIRECT_URI, SPOTIFY_CLIENT_ID } from "../../env"

const AUTH_URL = "https://accounts.spotify.com/authorize?"
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
        <code className="text-sm">
          accessToken: {localStorage.getItem("access_token")}
        </code>
        <code className="text-sm">
          refreshToken: {localStorage.getItem("refresh_token")}
        </code>
      </div>
    </main>
  )
}
