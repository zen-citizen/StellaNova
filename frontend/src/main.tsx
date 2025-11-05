import { APIProvider } from "@vis.gl/react-google-maps"
import { StrictMode } from "react"
import { createRoot } from "react-dom/client"
import App from "./App.tsx"
import { AppProvider } from "./contexts/AppContext.tsx"
import "./styles.css"

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <APIProvider apiKey={import.meta.env.VITE_GMAPS_API_KEY}>
      <AppProvider>
        <App />
      </AppProvider>
    </APIProvider>
  </StrictMode>
)
