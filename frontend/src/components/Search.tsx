import { useCallback, useState } from "react"
import useAppContext from "../hooks/useAppContext"
import useAutocompleteSuggestions from "../hooks/useAutocompleteSuggestions"

function Search() {
  const { setLocation } = useAppContext()

  const [inputValue, setInputValue] = useState<string>("")
  const { suggestions, resetSession } = useAutocompleteSuggestions(inputValue, {
    locationRestriction: {
      north: 13.5,
      south: 12.5,
      east: 78.25,
      west: 77
    }
  })

  const handleInput = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setInputValue(event.target.value)
    },
    [setInputValue]
  )

  const handleSelect = useCallback(
    (suggestion: google.maps.places.AutocompleteSuggestion) => {
      const place = suggestion.placePrediction!.toPlace()
      place
        .fetchFields({
          fields: ["location"]
        })
        .then(() => {
          const lat = place.location?.lat().toFixed(5)
          const lng = place.location?.lng().toFixed(5)
          if (lat && lng) {
            setLocation(Number(lat), Number(lng))
          }

          resetSession()
          setInputValue("")
        })
    },
    [resetSession, setLocation]
  )

  return (
    <div className="absolute top-2 left-2 w-full z-50">
      <input
        value={inputValue}
        onInput={handleInput}
        placeholder="Search for a place"
      />

      {suggestions.length > 0 && (
        <ul className="custom-list">
          {suggestions.map((suggestion, index) => {
            return (
              <li
                key={index}
                className="custom-list-item"
                onClick={() => handleSelect(suggestion)}
              >
                {suggestion.placePrediction!.text.text}
              </li>
            )
          })}
        </ul>
      )}
    </div>
  )
}

export default Search
