import Map from "./components/Map"
import Search from "./components/Search"
import Sidebar from "./components/Sidebar"

function App() {
  return (
    <div>
      <div className="grid md:h-screen grid-rows-[1fr_auto] md:grid-cols-[480px_1fr] md:grid-rows-1">
        <div className="hidden md:block">
          <Sidebar />
        </div>

        <div className="relative">
          <Search />
          <Map />
        </div>
      </div>

      <div className="w-full md:hidden">
        <Sidebar />
      </div>
    </div>
  )
}

export default App
