import { ChevronUp } from "lucide-react"
import type { ReactNode } from "react"

type AccordionProps = {
  title: string
  children: ReactNode
  defaultOpen?: boolean
}

function Accordion({ title, children, defaultOpen = false }: AccordionProps) {
  return (
    <details className="group" open={defaultOpen}>
      <summary className="flex cursor-pointer items-center justify-between py-1 px-0">
        <h3 className="text-heading-xs font-medium text-primary">{title}</h3>
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
