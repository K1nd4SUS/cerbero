import { Button, Spinner, Tab, Tabs, Tooltip } from "@nextui-org/react"
import { FaPen, FaTrash, FaArrowLeft, FaArrowRight } from "react-icons/fa6"
import { useFetchSync } from "../hooks/useFetch"
import { CerberoRegexes } from "../types/cerbero"

export type ServiceRegexesListProps = {
  nfq: string
}

export default function ServiceRegexesList({ nfq }: ServiceRegexesListProps) {
  const [
    response,
    isLoading,
    error
  ] = useFetchSync<CerberoRegexes>(`/api/regexes/${nfq}`)

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
      <Tab key="active" title="Active" className="flex flex-col h-full">
        {response?.regexes.active.length === 0 ?
          <div className="h-full w-full flex flex-col items-center justify-center">
            <span className="font-bold text-zinc-600">No regexes here, cerbero went to bed...</span>
          </div> :
          <ul className="h-full w-full flex flex-col gap-1">
            {response?.regexes.active.map((regex, i) => {
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
            <span className="font-bold text-zinc-600">No regexes here, cerbero went to bed...</span>
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
