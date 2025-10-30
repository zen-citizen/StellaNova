import maplibregl from "maplibre-gl"
import "maplibre-gl/dist/maplibre-gl.css"
import { useEffect, useRef } from "react"
import useAppContext from "../hooks/useAppContext"

function MapContainer() {
  const mapContainerRef = useRef<HTMLDivElement | null>(null)
  const mapRef = useRef<maplibregl.Map | null>(null)
  const markerRef = useRef<maplibregl.Marker | null>(null)
  const { setLocation, location } = useAppContext()

  useEffect(() => {
    if (!mapContainerRef.current) return

    const map = new maplibregl.Map({
      container: mapContainerRef.current,
      style: "https://basemaps.cartocdn.com/gl/positron-gl-style/style.json",
      minZoom: 9,
      maxZoom: 16,
      maxBounds: [
        [77, 12.5],
        [78.25, 13.5]
      ],
      center: [77.6, 12.974],
      zoom: 9.6
    })

    mapRef.current = map

    map.on("click", (e) => {
      setLocation(
        Number(e.lngLat.lat.toFixed(5)),
        Number(e.lngLat.lng.toFixed(5))
      )
    })

    return () => {
      markerRef.current?.remove()
      markerRef.current = null
      mapRef.current?.remove()
      mapRef.current = null
    }
  }, [setLocation])

  useEffect(() => {
    const map = mapRef.current
    if (!map) return

    if (location) {
      if (!markerRef.current) {
        markerRef.current = new maplibregl.Marker()
          .setLngLat([location.lng, location.lat])
          .addTo(map)
      } else {
        markerRef.current.setLngLat([location.lng, location.lat])
      }

      map.flyTo({
        center: [location.lng, location.lat],
        zoom: 12,
        speed: 1.5,
        curve: 1.6,
        easing: (t) => t
      })
    } else {
      markerRef.current?.remove()
      markerRef.current = null
    }
  }, [location])

  return <div ref={mapContainerRef} className="h-[50dvh] w-full md:h-full " />
}

export default MapContainer
