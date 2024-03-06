import { Button, Input, Spinner, Tab, Tabs, Tooltip } from "@nextui-org/react"
import { useEffect, useState } from "react"
import { FaArrowLeft, FaArrowRight, FaCircleCheck, FaPen, FaTrash } from "react-icons/fa6"
import { useFetch, useFetchSync } from "../hooks/useFetch"
import { CerberoRegexes } from "../types/cerbero"

export type ServiceRegexesListProps = {
  nfq: string
}

type NewRegexResponse = {
  regexes: {
    active: string[]
    inactive: string[]
  }
}

export default function ServiceRegexesList({ nfq }: ServiceRegexesListProps) {
  const [
    response,
    fetchRegexes,
    isLoading,
    error
  ] = useFetch<CerberoRegexes>()

  const [newRegex, setNewRegex] = useState("")

  const [
    ,
    newRegexFetch,
    isNewRegexLoading,
  ] = useFetch<NewRegexResponse>()

  useEffect(() => {
    fetchRegexes(`/api/regexes/${nfq}`)
  }, [])

  async function addNewRegex() {
    if(!newRegex) {
      return
    }

    await newRegexFetch(`/api/regexes/${nfq}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        regexes: [newRegex]
      })
    })

    setNewRegex("")
    fetchRegexes(`/api/regexes/${nfq}`)
  }

  if(isLoading) {
    return (
      <div className="h-full w-full flex flex-col items-center justify-center">
        <Spinner/>
      </div>
    )
  }

  if(error) {
    return (
      <div className="h-full w-full flex flex-col items-center justify-center">
        <span className="font-black text-xl text-zinc-300">{error.status} ({error.statusText})</span>
        <span className="font-semibold text-zinc-600">{error.data.error}</span>
      </div>
    )
  }

  return (
    <Tabs aria-label="Options">
      <Tab key="active" title="Active" className="h-full flex flex-col gap-2">
        <div className="flex items-center gap-2 p-2 rounded-lg bg-default-200">
          <Input
            placeholder="New regex"
            type="text"
            variant="flat"
            value={newRegex}
            onChange={({ target }) => setNewRegex(target.value)}
            className="bg-transparent"
          />
          <Button isLoading={isNewRegexLoading} variant="flat" color="success" onPress={addNewRegex}>
            <span className="font-bold">Add regex</span>
          </Button>
        </div>
        {response?.regexes.active.length === 0 ?
          <div className="h-full w-full flex flex-col items-center justify-center">
            <span className="font-bold text-zinc-600">No regexes here, add one from the input field above!</span>
          </div> :
          <ul className="h-full w-full flex flex-col gap-1">
            {response?.regexes.active.map((regex, i) => {
              return (
                <li key={i} className="text-sm bg-default-200 px-4 py-2 rounded-lg hover:opacity-75">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <span className="h-2 w-2 rounded-full bg-success"></span>
                      <span className="font-mono overflow-hidden line-clamp-1">{regex}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Tooltip content="Edit regex" delay={1000} size="sm">
                        <Button isIconOnly={true} color="warning" variant="flat" size="sm">
                          <FaPen/>
                        </Button>
                      </Tooltip>
                      <Tooltip content="Deactivate regex" delay={1000} size="sm">
                        <Button isIconOnly={true} color="danger" variant="flat" size="sm">
                          <FaArrowRight/>
                        </Button>
                      </Tooltip>
                      <Tooltip content="Delete regex" delay={1000} size="sm">
                        <Button isIconOnly={true} color="danger" variant="flat" size="sm">
                          <FaTrash/>
                        </Button>
                      </Tooltip>
                    </div>
                  </div>
                </li>
              )
            })}
          </ul>}
      </Tab>
      <Tab key="inactive" title="Inactive" className="flex flex-col h-full">
        {response?.regexes.inactive.length === 0 ?
          <div className="h-full w-full flex flex-col items-center justify-center">
            <span className="font-bold text-zinc-600">No regexes here, deactivate one from the `Active` tab.</span>
          </div> :
          <ul className="h-full w-full flex flex-col gap-1">
            {response?.regexes.inactive.map((regex, i) => {
              return (
                <li key={i} className="text-sm bg-default-200 px-4 py-2 rounded-lg hover:opacity-75">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <span className="h-2 w-2 rounded-full bg-success"></span>
                      <span className="font-mono">{regex}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Tooltip content="Edit regex" delay={1000} size="sm">
                        <Button isIconOnly={true} color="warning" variant="flat" size="sm">
                          <FaPen/>
                        </Button>
                      </Tooltip>
                      <Tooltip content="Deactivate regex" delay={1000} size="sm">
                        <Button isIconOnly={true} color="success" variant="flat" size="sm">
                          <FaArrowLeft/>
                        </Button>
                      </Tooltip>
                      <Tooltip content="Delete regex" delay={1000} size="sm">
                        <Button isIconOnly={true} color="danger" variant="flat" size="sm">
                          <FaTrash/>
                        </Button>
                      </Tooltip>
                    </div>
                  </div>
                </li>
              )
            })}
          </ul>}
      </Tab>
    </Tabs>
  )
}
