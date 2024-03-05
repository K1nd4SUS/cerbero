export function isNameInvalid(name: string) {
  return name.length === 0 || name.length > 32
}

export function isNfqInvalid(nfq: number) {
  return nfq < 100 || nfq > 199
}

export function isPortInvalid(port: number) {
  return port < 0 || port > 65535
}
