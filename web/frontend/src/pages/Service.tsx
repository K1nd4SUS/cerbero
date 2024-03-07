import { Card, CardBody, Chip } from "@nextui-org/react"
import { useParams } from "react-router-dom"
import Header from "../components/Header"
import ServiceRegexesList from "../components/ServiceRegexesList"
import { useFetchSync } from "../hooks/useFetch"
import Main from "../layouts/Main"
import Page from "../layouts/Page"
import type { CerberoService } from "../types/cerbero"
import Error from "./Error"
import Loading from "./Loading"

export default function Service() {
  const { nfq } = useParams()
  const [
    service,
    isLoading,
    error
  ] = useFetchSync<CerberoService>(`/api/services/${nfq}`)

  if(error) {
    return (
      <Error error={error}/>
    )
  }

  if(isLoading) {
    return (
      <Loading
        text="Loading service"
      />
    )
  }

  return (
    <Page>
      <Header/>
      <Main>
        <div className="h-full w-full flex flex-col">
          <div className="flex justify-between p-8">
            <div className="flex flex-col gap-2 p-8 border-b-2 border-zinc-700">
              <div className="flex flex-col gap-4 md:flex-row md:items-center">
                <span className="font-black text-5xl">{service?.name}</span>
                <Chip variant="flat" color="success" className="font-bold text-lg">
                  <span className="font-mono">{service?.protocol}://vm:{service?.port}</span>
                </Chip>
              </div>
              <span className="font-mono text-3xl text-zinc-300">nfq:{service?.nfq}</span>
            </div>
          </div>
          <div className="h-full flex flex-col gap-4 p-4 md:flex-row">
            <div className="flex-1 flex flex-col gap-4 p-4">
              <span className="font-bold text-3xl text-zinc-300">Metrics</span>
              <div className="h-full w-full flex flex-col items-center justify-center bg-default-100 rounded-xl">
                <span className="font-thin text-xl italic">Placeholder</span>
              </div>
            </div>
            <div className="flex-1 flex flex-col gap-4 p-4">
              <span className="font-bold text-3xl text-zinc-300">Regexes</span>
              <Card className="h-full bg-default-100">
                <CardBody>
                  <ServiceRegexesList nfq={nfq}/>
                </CardBody>
              </Card>
            </div>
          </div>
        </div>
      </Main>
    </Page>
  )
}
