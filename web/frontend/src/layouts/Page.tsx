import { ReactNode } from "react"
import backgroundImg from "../assets/images/background.png"

export type PageProps = {
  children: ReactNode
}

export default function Page({ children }: PageProps) {
  return (
    <div style={{ backgroundImage: `url(${backgroundImg})` }} className="absolute h-full w-full flex flex-col bg-cover">
      {children}
    </div>
  )
}
