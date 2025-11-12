import { useState } from "react"
import useAppContext from "../hooks/useAppContext"
import Accordion from "./Accordion"
import Table from "./Table"

function Details() {
  const { loading, error, entities, address, setView } = useAppContext()
  const [openAccordions, setOpenAccordions] = useState<Record<string, boolean>>(
    {}
  )

  const handleToggle = (entityName: string, open: boolean) => {
    setOpenAccordions((prev) => ({
      ...prev,
      [entityName]: open
    }))
  }

  const handleGoBack = () => {
    setView("introduction")
  }

  return (
    <div className="flex flex-col gap-4 items-start">
      <button
        onClick={handleGoBack}
        className="text-secondary underline underline-offset-3 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer"
      >
        ← Go back
      </button>

      {loading && <div>Loading...</div>}
      {error && <div className="text-error">{error}</div>}
      {address && (
        <div className="flex flex-col gap-2">
          <h2 className="text-heading-s font-semibold text-primary">Address</h2>
          <p>{address}</p>
        </div>
      )}
      {entities && (
        <div className="divide-y divide-separator w-full">
          {entities.map((entity, index) => (
            <Accordion
              key={entity.name}
              title={entity.name}
              open={openAccordions[entity.name] ?? index === 0}
              onToggle={(open) => {
                handleToggle(entity.name, open)
              }}
            >
              {entity.is_available ? (
                <Table attributes={entity.attributes} />
              ) : (
                <div className="py-2">{entity.not_available_message}</div>
              )}
            </Accordion>
          ))}
        </div>
      )}
      {!loading && !error && !entities && <div>No data available</div>}
    </div>
  )
}

export default Details
