import {
  Combobox,
  ComboboxInput,
  ComboboxOption,
  ComboboxOptions
} from "@headlessui/react"
import { Search as SearchIcon, X } from "lucide-react"
import { useCallback, useState } from "react"
import useAppContext from "../hooks/useAppContext"
import useAutocompleteSuggestions from "../hooks/useAutocompleteSuggestions"

type AutocompleteSuggestion = google.maps.places.AutocompleteSuggestion

function Search() {
  const { setLocation } = useAppContext()

  const [inputValue, setInputValue] = useState<string>("")
  const [selectedSuggestion, setSelectedSuggestion] =
    useState<AutocompleteSuggestion | null>(null)
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
      setSelectedSuggestion(null)
    },
    []
  )

  const handleSelect = useCallback(
    (suggestion: AutocompleteSuggestion | null) => {
      if (!suggestion) return

      setInputValue(suggestion.placePrediction!.text.text)
      setSelectedSuggestion(suggestion)

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
        })
    },
    [resetSession, setLocation]
  )

  const handleClear = useCallback(() => {
    setInputValue("")
    setSelectedSuggestion(null)
    resetSession()
  }, [resetSession])

  return (
    <div className="absolute top-4 left-6 right-18 z-50">
      <Combobox value={selectedSuggestion} onChange={handleSelect}>
        <div className="relative">
          <div className="absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none">
            <SearchIcon className="w-5 h-5 text-accent" />
          </div>
          <ComboboxInput
            aria-label="Enter the exact address in Bengaluru"
            autoComplete="off"
            displayValue={(suggestion: AutocompleteSuggestion | null) =>
              suggestion ? suggestion.placePrediction!.text.text : inputValue
            }
            onChange={handleInput}
            className="w-full pl-10 pr-10 py-2.5 border border-separator rounded-lg bg-surface text-primary placeholder:text-muted/70 focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent shadow-sm"
            placeholder="Enter the exact address in Bengaluru"
          />
          {inputValue && (
            <button
              type="button"
              onClick={handleClear}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-primary focus:outline-none focus:text-primary cursor-pointer focus:ring-2 focus:ring-accent rounded-sm"
              aria-label="Clear search"
            >
              <X className="w-5 h-5" />
            </button>
          )}
        </div>
        <ComboboxOptions
          anchor="bottom"
          className="w-(--input-width) mt-2 max-h-60 overflow-auto rounded-lg border border-separator bg-surface shadow-sm empty:invisible z-50"
          modal={false}
        >
          {suggestions.map((suggestion, index) => (
            <ComboboxOption
              key={index}
              value={suggestion}
              className="px-4 py-2 cursor-pointer data-focus:bg-accent-lighter text-primary"
            >
              {suggestion.placePrediction!.text.text}
            </ComboboxOption>
          ))}
        </ComboboxOptions>
      </Combobox>
    </div>
  )
}

export default Search
