export type Address = {
  text: string
  link: string
}

export type Attribute = {
  name: string
  value: string
  address?: Address
}

export type Entity = {
  name: string
  is_available: boolean
  not_available_message: string
  disclaimer?: string
  attributes: Attribute[]
}

export type EntitiesResponse = {
  entities: Entity[]
}

export type EntitiesError = {
  error: string
}

export type Location = {
  lat: number
  lng: number
}

export type View = "introduction" | "details"

export type Source = "map" | "search"

export type AddressResponse = {
  results: {
    formatted_address: string
  }[]
  status: string
}
