import { Button } from "@nextui-org/react"
import { useState } from "react"
import cerberoPng from "../assets/images/cerbero.png"
import SetupDialog from "../components/SetupDialog"
import Page from "../layouts/Page"
import { CerberoServiceInput } from "../types/cerbero"

export default function Home() {
  const [isSetupDialogOpen, setIsSetupDialogOpen] = useState(false)
  const [services, setServices] = useState<CerberoServiceInput[]>([{
    name: "",
    nfq: "",
    port: "",
    protocol: "tcp"
  }])

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
          <Button variant="flat" color="success" onPress={() => setIsSetupDialogOpen(true)}>
            <span className="font-bold">Setup Cerbero</span>
          </Button>
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
