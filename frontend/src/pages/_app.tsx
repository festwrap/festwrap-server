import type { AppProps } from "next/app"
import { Poppins } from "next/font/google"
import "../styles/globals.css"
import SessionWrapper from "@/components/SessionWrapper"
import RootLayout from "@/components/RootLayout"
import { twMerge } from "tailwind-merge"

const poppins = Poppins({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
})

export default function MyApp({ Component, pageProps }: AppProps) {
  return (
    <SessionWrapper>
      <main className={poppins.className}>
        <RootLayout>
          <Component {...pageProps} />
        </RootLayout>
      </main>
    </SessionWrapper>
  )
}
