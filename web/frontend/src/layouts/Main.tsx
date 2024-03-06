import type { ReactNode } from "react"
import Sidebar from "../components/Sidebar"

export type MainProps = {
  children: ReactNode
}

export default function Main({ children }: MainProps) {
  return (
    <main className="h-full w-full flex overflow-auto">
      <Sidebar/>
      {children}
    </main>
  )
}
