import { Card, CardBody, CardHeader, Chip, Divider, Skeleton } from "@nextui-org/react"
import { Link } from "react-router-dom"
import { FaCode } from "react-icons/fa6"
import { useFetchSync } from "../hooks/useFetch"
import type { CerberoRegexes } from "../types/cerbero"

export type ServiceProps = {
  chain: string
  name: string
  nfq: number
  port: number
  protocol: "tcp" | "udp"
}

export default function ServiceCard({ chain, name, nfq, port, protocol }: ServiceProps) {
  const [
    regexes,
    isLoading
  ] = useFetchSync<CerberoRegexes>(`/api/regexes/${nfq}`)

  return (
    <Link to={`/services/${nfq}`}>
      <Card key={nfq} className="h-full w-full text-zinc-300 hover:scale-[102.5%] hover:cursor-pointer">
        <CardHeader className="w-full flex items-center">
          <div className="flex flex-col gap-1 px-2">
            <span className="flex items-center gap-2">
              <FaCode className="text-2xl"/>
              <span className="font-black text-md">{name}</span>
            </span>
            <span className="font-mono text-sm text-default-500">nfq:{nfq}</span>
          </div>
          <div className="flex items-center gap-2 ml-auto px-2">
            <Chip variant="flat" color="success">
              <span className="font-mono">{protocol}://vm:{port}</span>
            </Chip>
            <Chip variant="flat" color="primary">
              <span className="font-mono">chain:{chain}</span>
            </Chip>
          </div>
        </CardHeader>
        <Divider/>
        <CardBody>
          <div className="h-full flex items-center justify-evenly">
            <Skeleton isLoaded={!isLoading} className="rounded-lg">
              <span className="text-green-600">
                {regexes?.regexes.active.length ?? 0} active regexes
              </span>
            </Skeleton>
            <Skeleton isLoaded={!isLoading} className="rounded-lg">
              <span className="text-zinc-600">
                {regexes?.regexes.inactive.length ?? 0} inactive regexes
              </span>
            </Skeleton>
            <Skeleton isLoaded={!isLoading} className="rounded-lg">
              <span className="text-kinda-accent">
                0 dropped packets
              </span>
            </Skeleton>
          </div>
        </CardBody>
      </Card>
    </Link>
  )
}
