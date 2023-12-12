import type { CerberoRegexes } from "./regex"

export type CerberoService = {
  name: string
  nfq: number
  port: number
  protocol: "tcp" | "udp"
  regexes?: CerberoRegexes[]
}

export type CerberoServiceCreate = Omit<CerberoService, "regexes"> & {
  regexes?: string[]
}
