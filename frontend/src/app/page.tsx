"use client"
import { useSession, signIn, signOut } from "next-auth/react"
import Button from "@components/Button"
import Card from "@components/Card"
import { PUBLIC_SPOTIFY_REDIRECT_URI } from "../../env"

export default function Home() {
  const { data: session } = useSession()

  if (session) {
    return (
      <Card>
        <span>Signed in as {session?.user?.email || session?.user?.name}</span>
        <Button accent="tertiary" onClick={() => signOut()}>
          Sign out
        </Button>
      </Card>
    )
  }

  return (
    <Card>
      <span>Not signed in</span>
      <Button
        accent="secondary"
        onClick={() =>
          signIn("spotify", { callbackUrl: PUBLIC_SPOTIFY_REDIRECT_URI })
        }
      >
        Login with Spotify
      </Button>
    </Card>
  )
}
