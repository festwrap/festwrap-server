import Link from "next/link"

const NavLink = ({
  href,
  children,
}: {
  href: string
  children: React.ReactNode
}) => {
  return (
    <Link
      href={href}
      className="text-primary hover:text-light transition-all ease-in-out duration-200 font-medium"
    >
      {children}
    </Link>
  )
}

export default NavLink
