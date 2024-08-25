"use client"
import { useSession, signIn, signOut } from "next-auth/react"
import Button from "@components/Button"

export default function Home() {
  const { data: session } = useSession()

  if (session) {
    return (
      <>
        Signed in as {session?.user?.email} <br />
        <Button accent="tertiary" onClick={() => signOut()}>
          Sign out
        </Button>
      </>
    )
  }

  return (
    <>
      Not signed in <br />
      <Button
        accent="secondary"
        onClick={() =>
          signIn("spotify", { callbackUrl: "http://localhost:3000" })
        }
      >
        Login with Spotify
      </Button>
    </>
  )
}
