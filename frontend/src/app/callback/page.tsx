"use client"

import { useRouter, useSearchParams } from "next/navigation"
import querystring from "querystring"
import { useEffect } from "react"
import axios, { AxiosResponse } from "axios"
import { SPOTIFY_CLIENT_ID, SPOTIFY_SECRET } from "../../../env"

const TOKEN_URL = "https://accounts.spotify.com/api/token"

export async function requestSpotifyToken(
  authCode: string
): Promise<{ access_token: string; refresh_token: string }> {
  const tokenRequestOptions = {
    method: "post",
    url: TOKEN_URL,
    data: querystring.stringify({
      code: authCode,
      redirect_uri: process.env.NEXT_PUBLIC_REDIRECT_URI,
      grant_type: "authorization_code",
    }),
    headers: {
      "content-type": "application/x-www-form-urlencoded",
      Authorization:
        "Basic " +
        Buffer.from(SPOTIFY_CLIENT_ID + ":" + SPOTIFY_SECRET).toString(
          "base64"
        ),
    },
  }

  return axios(tokenRequestOptions)
    .then((authResponse: AxiosResponse<any>) => {
      return authResponse.data
    })
    .catch((error: any) => {
      throw new Error(`Could not get access token: ${error}`)
    })
}

const Callback = () => {
  const searchParams = useSearchParams()
  const router = useRouter()

  const code = searchParams.get("code")
  const state = searchParams.get("state")

  useEffect(() => {
    async function fetchData() {
      if (code) {
        const response = await requestSpotifyToken(code)

        localStorage.setItem("access_token", response.access_token)
        localStorage.setItem("refresh_token", response.refresh_token)
      }
    }

    fetchData()
    router.push("/")
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return <div>Redirecting...</div>
}

export default Callback
