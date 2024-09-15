import Image from "next/image"
import Link from "next/link"

const Header = () => {
  return (
    <header className="flex items-center justify-between text-dark px-10 h-28">
      <h1 className="text-xl font-bold">
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
      </h1>
      <nav>
        <ul className="flex gap-6 font-medium">
          <li>
            <Link href="/get-started">Get started</Link>
          </li>
          <li>
            <Link href="/how-it-works">How does it works?</Link>
          </li>
          <li>
            <Link href="/about-us">About us</Link>
          </li>
        </ul>
      </nav>
    </header>
  )
}

export default Header
