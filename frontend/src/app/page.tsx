"use client"
import { useSession, signIn, signOut } from "next-auth/react"

export default function Home() {
  const { data: session } = useSession()

  if (session) {
    return (
      <>
        Signed in as {session?.user?.email} <br />
        <button
          className="hover:bg-gray-100 border border-gray-200 px-4 py-2 rounded-md"
          onClick={() => signOut()}
        >
          Sign out
        </button>
      </>
    )
  }

  return (
    <>
      Not signed in <br />
      <button
        className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-md"
        onClick={() =>
          signIn("spotify", { callbackUrl: "http://localhost:3000" })
        }
      >
        Login with Spotify
      </button>
    </>
  )
}
