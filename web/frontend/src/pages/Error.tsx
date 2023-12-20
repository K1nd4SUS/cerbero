import Header from "../components/Header"
import type { UseFetchError } from "../hooks/useFetch"
import Main from "../layouts/Main"
import Page from "../layouts/Page"

export type ErrorProps = {
  error: UseFetchError
}

export default function Error({ error }: ErrorProps) {
  return (
    <Page>
      <Header/>
      <Main>
        <div className="h-full w-full flex flex-col items-center justify-center">
          <span className="font-black text-6xl text-zinc-300">Error {error.status} ({error.statusText})</span>
          <p className="font-semibold text-3xl text-zinc-600">{error.data.error}.</p>
        </div>
      </Main>
    </Page>
  )
}
