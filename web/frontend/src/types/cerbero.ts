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

export type CerberoRegexes = Required<Pick<CerberoService, "regexes">>
