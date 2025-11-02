import { createContext, type ReactNode, useCallback, useState } from "react"
import { fetchAddress, fetchEntities } from "../api"
import type { Entity, Location, View } from "../types"

type AppContextType = {
  location: Location | null
  entities: Entity[] | null
  address: string | null
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

function getErrorMessage(err: unknown): string {
  return err instanceof Error ? err.message : "An unknown error occurred"
}

function AppProvider({ children }: AppProviderProps) {
  const [location, setLocationState] = useState<Location | null>(null)
  const [entities, setEntities] = useState<Entity[] | null>(null)
  const [address, setAddress] = useState<string | null>(null)
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string | null>(null)
  const [view, setView] = useState<View>("introduction")

  const setLocation = useCallback(async (lat: number, lng: number) => {
    const newLocation = { lat, lng }
    setLocationState(newLocation)
    setView("details")
    setLoading(true)
    setError(null)
    setEntities(null)
    setAddress(null)

    try {
      const results = await Promise.allSettled([
        fetchEntities(lat, lng),
        fetchAddress(lat, lng)
      ])

      let firstError: string | null = null

      if (results[0].status === "fulfilled") {
        setEntities(results[0].value)
      } else {
        firstError = getErrorMessage(results[0].reason)
        setEntities(null)
      }

      if (results[1].status === "fulfilled") {
        setAddress(results[1].value)
      } else {
        const errorMessage = getErrorMessage(results[1].reason)
        if (!firstError) {
          firstError = errorMessage
        }
        setAddress(null)
      }

      if (firstError) {
        setError(firstError)
      } else {
        setError(null)
      }
    } catch (err) {
      setError(getErrorMessage(err))
      setEntities(null)
      setAddress(null)
    } finally {
      setLoading(false)
    }
  }, [])

  const value: AppContextType = {
    location,
    entities,
    address,
    loading,
    error,
    view,
    setLocation,
    setView
  }

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>
}

export { AppContext, AppProvider }
