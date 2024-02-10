import { Spinner } from "@nextui-org/react"
import Header from "../components/Header"
import Main from "../layouts/Main"
import Page from "../layouts/Page"

export type LoadingProps = {
  text?: string
}

export default function Loading({ text }: LoadingProps) {
  return (
    <Page>
      <Header/>
      <Main>
        <div className="h-full w-full flex flex-col items-center justify-center">
          <div className="flex flex-col gap-4">
            <span className="font-bold text-xl">{text}</span>
            <Spinner/>
          </div>
        </div>
      </Main>
    </Page>
  )
}
