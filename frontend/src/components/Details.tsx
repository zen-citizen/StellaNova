import useAppContext from "../hooks/useAppContext"
import Accordion from "./Accordion"
import Table from "./Table"

function Details() {
  const { loading, error, data, setView } = useAppContext()

  const handleGoBack = () => {
    setView("introduction")
  }

  return (
    <div>
      <div className="mb-4">
        <button
          onClick={handleGoBack}
          className="text-blue-500 underline hover:text-blue-700 transition-colors"
        >
          ← Go back
        </button>
      </div>

      {loading && <div>Loading...</div>}
      {error && <div className="text-red-500">{error}</div>}
      {data && data.entities && (
        <div className="divide-y divide-separator">
          {data.entities.map((entity, index) => (
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
      {!loading && !error && !data && <div>No data available</div>}
    </div>
  )
}

export default Details
