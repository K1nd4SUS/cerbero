export type CerberoService = {
  chain: string
  name: string
  nfq: number
  port: number
  protocol: "tcp" | "udp"
  regexes?: {
    active: string[]
    inactive: string[]
  }
}

export type CerberoServiceInput = {
  chain: string
  name: string
  nfq: string
  port: string
  protocol: "tcp" | "udp"
}

export type CerberoRegexes = Required<Pick<CerberoService, "regexes">>
