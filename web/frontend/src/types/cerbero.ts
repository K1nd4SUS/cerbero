export type CerberoService = {
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
  name: string
  nfq: string
  port: string
  protocol: "tcp" | "udp"
}

export type CerberoRegexes = Required<Pick<CerberoService, "regexes">>
