import { Button } from "@nextui-org/react"
import { useState } from "react"
import { useNavigate } from "react-router-dom"
import cerberoPng from "../assets/images/cerbero.png"
import SetupDialog from "../components/SetupDialog"
import { useFetchSync } from "../hooks/useFetch"
import Page from "../layouts/Page"
import Error from "../pages/Error"
import { CerberoServiceInput } from "../types/cerbero"

export default function Home() {
  const [
    response,
    isLoading,
    error
  ] = useFetchSync<{ isSetupDone: boolean }>("/api/setup")
  const [isSetupDialogOpen, setIsSetupDialogOpen] = useState(false)
  const [services, setServices] = useState<CerberoServiceInput[]>([{
    name: "",
    nfq: "",
    port: "",
    protocol: "tcp"
  }])

  const navigate = useNavigate()

  if(error) {
    return (
      <Error error={error}/>
    )
  }

  return (
    <Page>
      <div className="h-full flex flex-col gap-8 items-center justify-center">
        <div className="flex flex-col gap-2 items-center">
          <img
            src={cerberoPng}
            alt="cerbero"
            className="h-16"
          />
          <h2 className="font-thin text-xl text-kinda-accent">A packet filtering tool for A/D CTFs</h2>
        </div>
        <div className="flex items-center gap-4">
          <Button variant="flat" color="warning">
            <span className="font-bold">Documentation</span>
          </Button>
          {response?.isSetupDone ?
            <Button isLoading={isLoading} variant="flat" color="success" onPress={() => navigate("/services")}>
              <span className="font-bold">Go to services</span>
            </Button> :
            <Button isLoading={isLoading} variant="flat" color="success" onPress={() => setIsSetupDialogOpen(true)}>
              <span className="font-bold">Setup Cerbero</span>
            </Button>}
        </div>
      </div>
      <SetupDialog
        isOpen={isSetupDialogOpen}
        services={services}
        setIsOpen={setIsSetupDialogOpen}
        setServices={setServices}
      />
    </Page>
  )
}
