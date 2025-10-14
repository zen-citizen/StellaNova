import useAppContext from "../hooks/useAppContext"

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
      {data && (
        <pre className="whitespace-pre-wrap text-sm overflow-auto">
          {JSON.stringify(data, null, 2)}
        </pre>
      )}
      {!loading && !error && !data && <div>No data available</div>}
    </div>
  )
}

export default Details
