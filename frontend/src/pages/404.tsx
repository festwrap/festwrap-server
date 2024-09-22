import Button from "@/components/Button"
import Image from "next/image"
import Link from "next/link"

export default function Custom404() {
  return (
    <div className="flex flex-col gap-6 items-center justify-center">
      <Image
        src="/not-found.svg"
        alt="404"
        width="150"
        height="150"
        className="mx-auto"
      />
      <div className="flex flex-col gap-3 justify-center items-center">
        <h1 className="text-3xl text-light font-medium text-center">404</h1>
        <p className="text-light text-center">
          It seems like you got a little bit lost
        </p>
        <Link href="/" className="text-light">
          Go back home
        </Link>
      </div>
    </div>
  )
}
