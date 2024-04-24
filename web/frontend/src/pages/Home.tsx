import { Button } from "@nextui-org/react"
import { useState } from "react"
import { Link } from "react-router-dom"
import cerberoPng from "../assets/images/cerbero.png"
import SetupDialog from "../components/SetupDialog"
import { useFetchSync } from "../hooks/useFetch"
import { CerberoServiceInput } from "../types/cerbero"

export default function Home() {
  const [
    response,
    isLoading
  ] = useFetchSync<{ isSetupDone: boolean }>("/api/setup")

  const [isSetupDialogOpen, setIsSetupDialogOpen] = useState(false)
  const [services, setServices] = useState<CerberoServiceInput[]>([{
    chain: "",
    name: "",
    nfq: "",
    port: "",
    protocol: "tcp"
  }])

  return (
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
        <Button as="a" href="https://github.com/k1nd4sus/cerbero" target="_blank" variant="flat" color="warning">
          <span className="font-bold">Documentation</span>
        </Button>
        {response?.isSetupDone ?
          <Button as={Link} to="/services" isLoading={isLoading} variant="flat" color="success">
            <span className="font-bold">Go to services</span>
          </Button> :
          <Button isLoading={isLoading} variant="flat" color="success" onPress={() => setIsSetupDialogOpen(true)}>
            <span className="font-bold">Setup Cerbero</span>
          </Button>}
      </div>
      <SetupDialog
        isOpen={isSetupDialogOpen}
        services={services}
        setIsOpen={setIsSetupDialogOpen}
        setServices={setServices}
      />
    </div>
  )
}
