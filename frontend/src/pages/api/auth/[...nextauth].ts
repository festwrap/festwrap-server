import NextAuth from "next-auth"
import SpotifyProvider from "next-auth/providers/spotify"

const scope =
  "user-read-recently-played user-read-playback-state user-top-read user-modify-playback-state user-read-currently-playing user-follow-read playlist-read-private user-read-email playlist-modify-private playlist-modify-public playlist-read-private user-read-private user-library-read playlist-read-collaborative"

export const authOptions = {
  providers: [
    SpotifyProvider({
      clientId: process.env.SPOTIFY_CLIENT_ID as string,
      clientSecret: process.env.SPOTIFY_SECRET as string,
      authorization: {
        params: {
          scope,
        },
      },
    }),
  ],
  callbacks: {
    async jwt({ token, account }: any) {
      if (account) {
        token.id = account.id
        token.expires_at = account.expires_at
        token.accessToken = account.access_token
      }
      return token
    },
    async session({ session, token }: any) {
      session.user = token
      return session
    },
  },
}

export default NextAuth(authOptions)
