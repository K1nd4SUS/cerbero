import type { ReactNode } from "react"

export type MainProps = {
  children: ReactNode
}

export default function Main({ children }: MainProps) {
  return (
    <main className="h-full w-full overflow-auto">
      {children}
    </main>
  )
}
