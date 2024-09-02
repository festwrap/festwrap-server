"use client"
import { useSession, signIn, signOut } from "next-auth/react"
import Button from "@components/Button"
import Card from "@components/Card"

export default function Home() {
  const { data: session } = useSession()

  const copyToClipboardToken = () => {
    navigator.clipboard.writeText(session?.user?.accessToken || "")
  }

  if (session) {
    return (
      <Card>
        <span>Signed in as {session?.user?.email || session?.user?.name}</span>
        <div className="flex flex-row gap-2">
          <Button accent="primary" onClick={copyToClipboardToken}>
            Copy access token
          </Button>
          <Button accent="tertiary" onClick={() => signOut()}>
            Sign out
          </Button>
        </div>
      </Card>
    )
  }

  return (
    <Card>
      <span>Not signed in</span>
      <Button
        accent="secondary"
        onClick={() => signIn("spotify", { callbackUrl: "/" })}
      >
        Login with Spotify
      </Button>
    </Card>
  )
}
