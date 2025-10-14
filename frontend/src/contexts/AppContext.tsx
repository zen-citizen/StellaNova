import { createContext, type ReactNode, useCallback, useState } from "react"
import type { EntitiesResponse, ErrorResponse, Location, View } from "../types"

type AppContextType = {
  location: Location | null
  data: EntitiesResponse | null
  loading: boolean
  error: string | null
  view: View
  setLocation: (lat: number, lng: number) => void
  setView: (view: View) => void
}

const AppContext = createContext<AppContextType | undefined>(undefined)

type AppProviderProps = {
  children: ReactNode
}

function AppProvider({ children }: AppProviderProps) {
  const [location, setLocationState] = useState<Location | null>(null)
  const [data, setData] = useState<EntitiesResponse | null>(null)
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string | null>(null)
  const [view, setView] = useState<View>("introduction")

  const setLocation = useCallback(async (lat: number, lng: number) => {
    const newLocation = { lat, lng }
    setLocationState(newLocation)
    setView("details")
    setLoading(true)
    setError(null)
    setData(null)

    try {
      const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
      const response = await fetch(
        `${apiBaseUrl}/api/v1/entities?lat=${lat}&lng=${lng}&city=bengaluru`
      )

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const result = await response.json()

      // Check if the response contains an error
      if ("error" in result) {
        const errorResponse = result as ErrorResponse
        setError(errorResponse.error)
        setData(null)
      } else {
        const entitiesResponse = result as EntitiesResponse
        setData(entitiesResponse)
        setError(null)
      }
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "An unknown error occurred"
      setError(errorMessage)
      setData(null)
    } finally {
      setLoading(false)
    }
  }, [])

  const value: AppContextType = {
    location,
    data,
    loading,
    error,
    view,
    setLocation,
    setView
  }

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>
}

export { AppContext, AppProvider }
