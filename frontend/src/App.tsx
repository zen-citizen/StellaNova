import Map from "./components/Map"
import Search from "./components/Search"
import Sidebar from "./components/Sidebar"

function App() {
  return (
    <div className="grid md:h-screen md:grid-cols-[480px_1fr]">
      <div className="order-1 md:order-0">
        <Sidebar />
      </div>

      <div className="relative">
        <Search />
        <Map />
      </div>
    </div>
  )
}

export default App
