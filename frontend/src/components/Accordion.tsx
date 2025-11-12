import { ChevronUp } from "lucide-react"
import type { ReactNode } from "react"
import { useEffect, useState } from "react"

type AccordionProps = {
  title: string
  children: ReactNode
  open: boolean
  onToggle: (open: boolean) => void
}

function Accordion({ title, children, open, onToggle }: AccordionProps) {
  const [isMounted, setMount] = useState(false)

  useEffect(() => {
    setMount(true)
  }, [])

  return (
    <details
      className="group py-2"
      open={open}
      onToggle={() => isMounted && onToggle(!open)}
    >
      <summary className="flex cursor-pointer items-center justify-between py-1 px-0">
        <h3 className="text-heading-s font-semibold text-primary">{title}</h3>
        <ChevronUp
          size={24}
          className="text-primary transition-transform duration-300 ease-in-out group-open:rotate-180"
        />
      </summary>
      {children}
    </details>
  )
}

export default Accordion
