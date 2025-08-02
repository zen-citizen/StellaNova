import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";

const elem = document.getElementById("root");

if (elem instanceof HTMLDivElement) {
  createRoot(elem).render(
    <StrictMode>
      <App />
    </StrictMode>
  );
}
