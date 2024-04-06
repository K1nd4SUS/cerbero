import {
  Button,
  Input,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Pagination,
  Radio,
  RadioGroup
} from "@nextui-org/react"
import { Dispatch, SetStateAction, useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { FaCheck, FaPlus, FaTrash } from "react-icons/fa6"
import { CerberoServiceInput } from "../types/cerbero"
import { isChainInvalid, isNameInvalid, isNfqInvalid, isPortInvalid } from "../utils/validation"

export type SetupDialogProps = {
  isOpen: boolean
  services: CerberoServiceInput[]
  setIsOpen: Dispatch<SetStateAction<boolean>>
  setServices: Dispatch<SetStateAction<CerberoServiceInput[]>>
}

export default function SetupDialog({ isOpen, services, setIsOpen, setServices }: SetupDialogProps) {
  const [currentServiceIndex, setCurrentServiceIndex] = useState(0)
  const [chain, setChain] = useState("")
  const [name, setName] = useState("")
  const [nfq, setNfq] = useState("")
  const [port, setPort] = useState("")
  const [protocol, setProtocol] = useState<"tcp" | "udp">("tcp")

  const navigate = useNavigate()

  // Add a new (empty) service
  function addService() {
    setChain("")
    setName("")
    setNfq("")
    setPort("")
    setProtocol("tcp")

    setServices(services.concat({
      chain: "",
      name: "",
      nfq: "",
      port: "",
      protocol: "tcp"
    }))

    setCurrentServiceIndex(services.length)
  }

  // Delete the currently selected service
  function deleteService() {
    if(!canServiceBeDeleted()) {
      return
    }

    const previousServiceIndex = currentServiceIndex

    setServices(services.filter((_, i) => i !== currentServiceIndex))
    setCurrentServiceIndex(previousServiceIndex - 1)
  }

  // Loads the service at the current index from the services array into the state variables
  function loadServiceIntoState() {
    if(!services[currentServiceIndex]) {
      return
    }

    setChain(services[currentServiceIndex].chain)
    setName(services[currentServiceIndex].name)
    setNfq(services[currentServiceIndex].nfq.toString())
    setPort(services[currentServiceIndex].port.toString())
    setProtocol(services[currentServiceIndex].protocol)
  }

  function canServiceBeDeleted() {
    return services.length > 1
  }

  function canSetupBeCompleted() {
    for(const service of services) {
      const integerNfq = parseInt(service.nfq)
      const integerPort = parseInt(service.port)

      if(isChainInvalid(service.chain)) {
        return false
      }

      if(isNameInvalid(service.name)) {
        return false
      }

      if(!parseInt(service.nfq) || isNfqInvalid(integerNfq)) {
        return false
      }

      if(!parseInt(service.port) || isPortInvalid(integerPort)) {
        return false
      }
    }

    return true
  }

  async function submitSetup() {
    await fetch("/api/services", {
      method: "POST",
      body: JSON.stringify(services.map(service => {
        return {
          ...service,
          nfq: parseInt(service.nfq),
          port: parseInt(service.port)
        }
      })),
      headers: {
        "Content-Type": "application/json"
      }
    })

    navigate("/services")
  }

  useEffect(() => {
    loadServiceIntoState()
  }, [currentServiceIndex])

  // Updates the services array (in the parent component) with the
  // new inserted data, this allows real-time modification of the services
  // (prevents clicking a save button to save the modifications on a service)
  useEffect(() => {
    const newServices = services.map((service, i) => {
      if(i === currentServiceIndex) {
        return {
          ...service,
          chain,
          name,
          nfq: nfq,
          port: port,
          protocol
        }
      }

      return service
    })

    setServices(newServices)
  }, [chain, name, nfq, port, protocol])

  return (
    <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} size="3xl">
      <ModalContent>
        <ModalHeader className="flex flex-col gap-1">Services Setup</ModalHeader>
        <ModalBody className="grid grid-cols-2">
          <Input
            label="Name"
            placeholder="Broken API"
            type="text"
            value={name}
            minLength={1}
            maxLength={32}
            onChange={({ target }) => setName(target.value)}
          />
          <Input
            label="Port"
            placeholder="1337"
            type="number"
            value={port}
            min={1}
            max={65535}
            isInvalid={isPortInvalid(parseInt(port))}
            onChange={({ target }) => setPort(target.value)}
          />
          <Input
            label="Nfq"
            placeholder="101"
            type="number"
            value={nfq}
            min={100}
            max={199}
            isInvalid={isNfqInvalid(parseInt(nfq))}
            onChange={({ target }) => setNfq(target.value)}
          />
          <Input
            label="Chain"
            placeholder="INPUT"
            type="text"
            value={chain}
            minLength={1}
            maxLength={32}
            onChange={({ target }) => setChain(target.value)}
          />
          <RadioGroup
            label="Protocol"
            size="sm"
            defaultChecked={true}
            defaultValue={protocol}
            value={protocol}
            onChange={({ target }) => setProtocol(target.value as "tcp" | "udp")}
            className="col-span-2 items-center p-2"
          >
            <div className="flex items-center gap-4">
              <Radio value="tcp">TCP</Radio>
              <Radio value="udp">UDP</Radio>
            </div>
          </RadioGroup>
        </ModalBody>
        <ModalFooter>
          <div className="w-full flex items-center justify-between">
            <Pagination total={services.length} page={currentServiceIndex + 1} onChange={pageNumber => setCurrentServiceIndex(pageNumber - 1)}/>
            <div className="flex items-center gap-2">
              {canSetupBeCompleted() ?
                <Button startContent={<FaCheck/>} color="success" variant="flat" size="sm" onPress={() => void submitSetup()}>
                  Complete setup
                </Button> : <></>}
              <Button startContent={<FaPlus/>} color="success" variant="flat" size="sm" onPress={addService}>
                Add service
              </Button>
              {canServiceBeDeleted() ?
                <Button startContent={<FaTrash/>} color="danger" variant="flat" size="sm" onPress={deleteService}>
                  Remove service
                </Button> : <></>}
            </div>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}
