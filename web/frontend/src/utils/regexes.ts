export function hexEncodeRegex(regex: string) {
  const encoder = new TextEncoder();
  const encodedArray = encoder.encode(regex);
  const encodedRegex = Array.from(encodedArray).map(byte => byte.toString(16).padStart(2, '0')).join('');

  return encodedRegex
}
