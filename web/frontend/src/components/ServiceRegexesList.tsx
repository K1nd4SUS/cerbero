import { Button, Input, Tab, Tabs, Tooltip } from "@nextui-org/react"
import { useEffect, useState } from "react"
import { FaArrowLeft, FaArrowRight, FaTrash } from "react-icons/fa6"
import Error from "../components/Error"
import Loading from "../components/Loading"
import { useFetch } from "../hooks/useFetch"
import { CerberoRegexes } from "../types/cerbero"
import { hexEncode } from "../utils/regexes"

export type ServiceRegexesListProps = {
  nfq: string
}

export default function ServiceRegexesList({ nfq }: ServiceRegexesListProps) {
  const [
    regexesResponse,
    regexesFetch,
    isRegexesLoading,
    regexesError
  ] = useFetch<CerberoRegexes>()

  const [newRegex, setNewRegex] = useState("")

  const [
    ,
    newRegexFetch,
    isNewRegexLoading,
  ] = useFetch()

  const [
    ,
    editRegexFetch,
    isEditRegexFetchLoading
  ] = useFetch()

  const [
    ,
    deleteRegexFetch,
    isDeleteRegexLoading,
  ] = useFetch()

  useEffect(() => {
    void regexesFetch(`/api/regexes/${nfq}`)
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
    void regexesFetch(`/api/regexes/${nfq}`)
  }

  async function editRegex(regex: string, currentState: "active" | "inactive", newRegex: string, newState: "active" | "inactive") {
    const reghex = hexEncode(regex)

    await editRegexFetch(`/api/regexes/${nfq}?reghex=${reghex}&state=${currentState}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        regex: newRegex,
        state: newState
      })
    })

    void regexesFetch(`/api/regexes/${nfq}`)
  }

  async function deleteRegex(regex: string, state: "active" | "inactive") {
    const reghex = hexEncode(regex)

    await deleteRegexFetch(`/api/regexes/${nfq}?reghex=${reghex}&state=${state}`, {
      method: "DELETE"
    })

    void regexesFetch(`/api/regexes/${nfq}`)
  }

  if(isRegexesLoading) {
    return (
      <Loading text="Loading regexes..."/>
    )
  }

  if(regexesError) {
    return (
      <Error error={regexesError}/>
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
          <Button isLoading={isNewRegexLoading} variant="flat" color="success" onPress={() => void addNewRegex()}>
            <span className="font-bold">Add regex</span>
          </Button>
        </div>
        {regexesResponse?.regexes.active.length === 0 ?
          <div className="h-full w-full flex flex-col items-center justify-center">
            <span className="font-bold text-zinc-600">No regexes here, add one from the input field above!</span>
          </div> :
          <ul className="h-full w-full flex flex-col gap-1">
            {regexesResponse?.regexes.active.map((regex, i) => {
              return (
                <li key={i} className="text-sm bg-default-200 px-4 py-2 rounded-lg hover:opacity-75">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <span className="h-2 w-2 rounded-full bg-success"></span>
                      <span className="font-mono overflow-hidden line-clamp-1">{regex}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Tooltip content="Deactivate regex" delay={1000} size="sm">
                        <Button isLoading={isEditRegexFetchLoading} onPress={() => void editRegex(regex, "active", regex, "inactive")} isIconOnly={true} color="danger" variant="flat" size="sm">
                          <FaArrowRight/>
                        </Button>
                      </Tooltip>
                      <Tooltip content="Delete regex" delay={1000} size="sm">
                        <Button isLoading={isDeleteRegexLoading} onPress={() => void deleteRegex(regex, "active")} isIconOnly={true} color="danger" variant="flat" size="sm">
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
        {regexesResponse?.regexes.inactive.length === 0 ?
          <div className="h-full w-full flex flex-col items-center justify-center">
            <span className="font-bold text-zinc-600">No regexes here, deactivate one from the `Active` tab.</span>
          </div> :
          <ul className="h-full w-full flex flex-col gap-1">
            {regexesResponse?.regexes.inactive.map((regex, i) => {
              return (
                <li key={i} className="text-sm bg-default-200 px-4 py-2 rounded-lg hover:opacity-75">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <span className="h-2 w-2 rounded-full bg-success"></span>
                      <span className="font-mono">{regex}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Tooltip content="Activate regex" delay={1000} size="sm">
                        <Button isLoading={isEditRegexFetchLoading} onPress={() => void editRegex(regex, "inactive", regex, "active")} isIconOnly={true} color="success" variant="flat" size="sm">
                          <FaArrowLeft/>
                        </Button>
                      </Tooltip>
                      <Tooltip content="Delete regex" delay={1000} size="sm">
                        <Button onPress={() => void deleteRegex(regex, "inactive")} isIconOnly={true} color="danger" variant="flat" size="sm">
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
