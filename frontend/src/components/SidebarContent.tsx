import useAppContext from "../hooks/useAppContext"
import Details from "./Details"
import Introduction from "./Introduction"

function SidebarContent() {
  const { view } = useAppContext()

  return (
    <div className="px-4 md:px-6 py-6">
      {view === "introduction" && <Introduction />}
      {view === "details" && <Details />}
    </div>
  )
}

export default SidebarContent
