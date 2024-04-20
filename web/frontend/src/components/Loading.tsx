import { Spinner } from "@nextui-org/react"

export type LoadingProps = {
  text?: string
}

export default function Loading({ text }: LoadingProps) {
  return (
    <div className="h-full w-full flex flex-col items-center justify-center">
      <div className="flex flex-col gap-4">
        <span className="font-bold text-xl">{text}</span>
        <Spinner/>
      </div>
    </div>
  )
}

