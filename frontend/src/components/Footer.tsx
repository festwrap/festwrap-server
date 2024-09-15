import Link from "next/link"

const Footer = () => {
  return (
    <footer className="flex flex-col items-center justify-center gap-3 bg-secondary bg-opacity-30 text-light px-10 h-36">
      <nav>
        <ul className="flex gap-6 font-medium text-sm">
          <li>
            <Link href="/get-started">Get started</Link>
          </li>
          <li>
            <Link href="/how-it-works">How does it works?</Link>
          </li>
          <li>
            <Link href="/about-us">About us</Link>
          </li>
          <li>
            <Link href="/terms-of-service">Terms of Service</Link>
          </li>
          <li>
            <Link href="/privacy-policy">Privacy Policy</Link>
          </li>
        </ul>
      </nav>
      <small>© 2021 Festwrap</small>
    </footer>
  )
}

export default Footer
