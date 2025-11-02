import useAppContext from "../hooks/useAppContext"
import Accordion from "./Accordion"
import Table from "./Table"

function Details() {
  const { loading, error, entities, address, setView } = useAppContext()

  const handleGoBack = () => {
    setView("introduction")
  }

  return (
    <div className="flex flex-col gap-4 items-start">
      <button
        onClick={handleGoBack}
        className="text-blue-500 underline hover:text-blue-700 transition-colors"
      >
        ← Go back
      </button>

      {loading && <div>Loading...</div>}
      {error && <div className="text-red-500">{error}</div>}
      {address && (
        <div className="flex flex-col gap-2">
          <h2 className="text-xl font-semibold">Address</h2>
          <p>{address}</p>
        </div>
      )}
      {entities && (
        <div className="divide-y divide-separator">
          {entities.map((entity, index) => (
            <Accordion
              key={entity.name}
              title={entity.name}
              defaultOpen={index === 0}
            >
              <Table attributes={entity.attributes} />
            </Accordion>
          ))}
        </div>
      )}
      {!loading && !error && !entities && <div>No data available</div>}
    </div>
  )
}

export default Details
