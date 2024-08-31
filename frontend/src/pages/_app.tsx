import type { AppProps } from "next/app"
import "../styles/globals.css"
import SessionWrapper from "@/components/SessionWrapper"

export default function MyApp({ Component, pageProps }: AppProps) {
  return (
    <SessionWrapper>
      <Component {...pageProps} />
    </SessionWrapper>
  )
}
