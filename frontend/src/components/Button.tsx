import { twMerge } from "tailwind-merge"

type AccentTypes = "primary" | "secondary" | "tertiary"

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode
  accent: AccentTypes
}

const accentClasses = {
  primary: "bg-primary text-white hover:bg-primary hover:bg-opacity-90",
  secondary: "bg-secondary text-dark hover:bg-opacity-90",
  tertiary:
    "bg-gray-100 text-gray-800 hover:bg-gray-100 border border-gray-100",
}

const baseClasses = "px-4 py-2 rounded-lg font-medium"

const Button = ({ children, accent, ...buttonProps }: ButtonProps) => {
  return (
    <button
      className={twMerge(baseClasses, accentClasses[accent])}
      {...buttonProps}
    >
      {children}
    </button>
  )
}

export default Button
