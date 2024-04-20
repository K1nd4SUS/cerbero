import { UseFetchError } from "../hooks/useFetch"

export type ErrorProps = {
  error: UseFetchError
}

export default function Error({ error }: ErrorProps) {
  return (
    <div className="h-full w-full flex flex-col items-center justify-center">
      <span className="font-black text-6xl text-zinc-300">Error {error.status} ({error.statusText})</span>
      <p className="font-semibold text-3xl text-zinc-600">{error.data.error}.</p>
    </div>
  )
}
