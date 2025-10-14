import useAppContext from "../hooks/useAppContext"

function MapContainer() {
  const { setLocation } = useAppContext()

  const handleSetLocation = () => {
    setLocation(12.97011, 77.64472)
  }

  return (
    <div className="grid h-full w-full place-items-center rounded-md">
      <div className="text-center">
        <div className="mb-4">MapContainer</div>
        <button
          onClick={handleSetLocation}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
        >
          Set Example Location
        </button>
      </div>
    </div>
  )
}

export default MapContainer
