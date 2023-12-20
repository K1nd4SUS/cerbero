import Header from "../components/Header"
import Service from "../components/Service"
import { useFetchSync } from "../hooks/useFetch"
import Main from "../layouts/Main"
import Page from "../layouts/Page"
import type { CerberoService } from "../types/cerbero"
import Error from "./Error"
import Loading from "./Loading"

export default function Services() {
  const [
    services,
    isLoading,
    error
  ] = useFetchSync<CerberoService[]>("/api/services")

  if(error) {
    return (
      <Error error={error}/>
    )
  }

  if(isLoading) {
    return (
      <Loading
        text="Loading services"
      />
    )
  }

  return (
    <Page>
      <Header/>
      <Main>
        <div className="h-full w-full flex flex-col gap-4 p-8">
          <span className="font-bold text-3xl text-zinc-300">Services</span>
          <div className="w-full grid grid-cols-1 gap-4 justify-center md:grid-cols-2">
            {services?.sort((a, b) => a.nfq - b.nfq).map(service => {
              return (
                <Service key={service.nfq} {...service}/>
              )
            })}
          </div>
        </div>
      </Main>
    </Page>
  )
}
