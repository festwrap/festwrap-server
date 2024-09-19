import Image from "next/image"
import Link from "next/link"

const Header = () => {
  return (
    <header className="text-dark">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between py-8">
          <div className="flex justify-start">
            <Link href="/">
              <Image
                src="/logo.svg"
                alt="Festwrap logo"
                width={150}
                height={150}
                className="h-auto"
                priority
              />
            </Link>
          </div>
          <nav className="hidden md:flex gap-10">
            <Link href="/get-started">Get started</Link>
            <Link href="/how-it-works">How does it works?</Link>
            <Link href="/about-us">About us</Link>
          </nav>
        </div>
      </div>
    </header>
  )
}

export default Header
