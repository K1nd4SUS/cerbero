export function hexEncode(s: string) {
  const encoder = new TextEncoder()
  const encodedArray = encoder.encode(s)
  const encodedS = Array.from(encodedArray).map(byte => byte.toString(16).padStart(2, "0")).join("")

  return encodedS
}
