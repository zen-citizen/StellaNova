import type { EntitiesError, EntitiesResponse, Entity } from "./types"

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

type Geocoder = google.maps.Geocoder

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

async function fetchAddress(
  lat: number,
  lng: number,
  geocoder: Geocoder | null
): Promise<string> {
  if (!geocoder) {
    throw new Error("Geocoder not initialized")
  }

  const response = await geocoder.geocode({
    location: {
      lat: lat,
      lng: lng
    }
  })

  if (response.results && response.results.length > 0) {
    return response.results[0].formatted_address
  } else {
    return "No address found"
  }
}

export { fetchAddress, fetchEntities }
