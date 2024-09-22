import type { Config } from "tailwindcss"

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#0D1B2A",
        secondary: "#E0E1DD",
        tertiary: "#778DA9",
        dark: "#0D1B2A",
        light: "#778DA9",
        spotify: "#00DA5A",
      },
    },
  },
  plugins: [],
}
export default config
