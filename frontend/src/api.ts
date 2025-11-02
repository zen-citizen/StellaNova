import type {
  AddressResponse,
  EntitiesError,
  EntitiesResponse,
  Entity
} from "./types"

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
const GMAPS_API_KEY = import.meta.env.VITE_GMAPS_API_KEY

async function fetchEntities(lat: number, lng: number): Promise<Entity[]> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/entities?lat=${lat}&lng=${lng}&city=bengaluru`
  )

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const result = await response.json()

  if ("error" in result) {
    const errorResponse = result as EntitiesError
    throw new Error(errorResponse.error)
  }

  const entitiesResponse = result as EntitiesResponse
  return entitiesResponse.entities
}

async function fetchAddress(lat: number, lng: number): Promise<string> {
  const response = await fetch(
    `https://maps.googleapis.com/maps/api/geocode/json?latlng=${lat},${lng}&key=${GMAPS_API_KEY}`
  )

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const result = (await response.json()) as AddressResponse

  if (result.results && result.results.length > 0) {
    return result.results[0].formatted_address
  } else {
    return "No address found"
  }
}

export { fetchAddress, fetchEntities }
