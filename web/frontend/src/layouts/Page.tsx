import { ReactNode } from "react"

export type PageProps = {
  children: ReactNode
}

export default function Page({ children }: PageProps) {
  return (
    <div className="absolute h-full w-full flex flex-col bg-default-50">
      {children}
    </div>
  )
}
