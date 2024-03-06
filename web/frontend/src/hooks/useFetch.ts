import { useEffect, useState } from "react"

export type UseFetchError = {
  status: number
  statusText: string
  data: {
    error: string
  }
}

/**
 * React hook abstraction for the fetch API.
 *
 * Returns an array containing 4 elements:
 *
 * 0. The `T` typed response if the request was successful, otherwise `undefined`.
 * 1. An async function that triggers the HTTP request whenever called.
 * 2. A `boolean` flag that states if the request is loading.
 * 3. In case of an error a `UseFetchError` object, otherwise `undefined`.
 */
export function useFetch<T>(): [
  T | undefined,
  (url: string, init?: RequestInit) => Promise<void>,
  boolean,
  UseFetchError | undefined
] {
  const [response, setResponse] = useState<T>()
  const [isLoading, setIsLoading] = useState(false)
  // TODO: set 400+ responses in response instead of error
  const [error, setError] = useState<UseFetchError>()

  async function triggerFetch(url: string, init?: RequestInit) {
    try {
      setIsLoading(true)

      const response = await fetch(url, init)
      const responseJson = await response.json() as T

      if(response.ok) {
        setResponse(responseJson)
      }
      else {
        setError({
          status: response.status,
          statusText: response.statusText,
          data: responseJson as UseFetchError["data"]
        })
      }
    }
    catch(e) {
      console.error(e)
    }
    finally {
      setIsLoading(false)
    }
  }

  return [
    response,
    triggerFetch,
    isLoading,
    error
  ]
}

/**
 * "Synchronous" version of the `useFetch` hook.
 * (Synchronous meaning that the HTTP request is triggered as soon as the hook is used).
 *
 * Returns an array containing 3 elements:
 *
 * 0. The `T` typed response if the request was successful, otherwise `undefined`.
 * 1. A `boolean` flag that states if the request is loading.
 * 2. In case of an error a `UseFetchError` object, otherwise `undefined`.
 */
export function useFetchSync<T>(url: string, init?: RequestInit): [
  T | undefined,
  boolean,
  UseFetchError | undefined
] {
  const [
    response,
    triggerFetch,
    isLoading,
    error
  ] = useFetch<T>()

  useEffect(() => {
    void triggerFetch(url, init)
  }, [])

  return [
    response,
    isLoading,
    error
  ]
}
