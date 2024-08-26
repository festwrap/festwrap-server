import Image from "next/image"

interface CardProps {
  children: React.ReactNode
}

const Card = ({ children }: CardProps) => {
  return (
    <div className="rounded-lg border border-gray-200 p-6 w-1/2 lg:w-1/3 gap-4 flex flex-col items-center justify-center">
      <Image
        src="/logo.svg"
        alt="Festwrap logo"
        width={150}
        height={150}
        className="h-auto"
        priority
      />
      {children}
    </div>
  )
}

export default Card
