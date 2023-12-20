import { nextui } from "@nextui-org/react"

/** @type {import("tailwindcss").Config} */
export default {
  content: [
    "**/*.{html,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        "kinda-accent": "#1CF1FB",
        "kinda-primary": "#141A26",
        "kinda-secondary": "#14222E"
      },
      fontFamily: {
        sans: ["manrope"],
      }
    }
  },
  plugins: [
    nextui({
      defaultTheme: "dark"
    })
  ]
}
