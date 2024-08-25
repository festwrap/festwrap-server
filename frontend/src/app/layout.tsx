import type { Metadata } from "next"
import { Poppins } from "next/font/google"
import "./globals.css"
import SessionWrapper from "./components/SessionWrapper"

const poppins = Poppins({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
})

export const metadata: Metadata = {
  title: "Festwrap",
  description: "Spotify playlist generator for music festivals",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <SessionWrapper>
      <html lang="en">
        <body className={poppins.className}>{children}</body>
      </html>
    </SessionWrapper>
  )
}
