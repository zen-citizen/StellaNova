import { Sidebar, type DataSource } from "./components";

const dataSources: DataSource[] = [
  {
    href: "https://opencity.in/data",
    displayName: "OpenCity",
  },
  {
    href: "https://kgis.ksrsac.in/kgis/",
    displayName: "Karnataka-GIS",
  },
  {
    href: "https://www.openstreetmap.org/about",
    displayName: "OpenStreetMap",
  },
];

const footerLinks: DataSource[] = [
  {
    href: "https://forms.gle/EmQiMpayciLdbww96",
    displayName: "Report an Error",
  },
  {
    href: "https://docs.google.com/forms/d/e/1FAIpQLScQS_-VgUFQZJedyu6iIlpoYymsKSyGUhrvPoJX1WkZGQqfLQ/viewform",
    displayName: "Volunteer with Us",
  },
  {
    href: "https://github.com/zen-citizen/StellaNova",
    displayName: "Open Source",
  },
];

function App() {
  return (
    <div className="relative isolate h-full min-h-screen">
      {/* Sidebar  */}
      <Sidebar {...{ footerLinks, dataSources }} />
      <div className=""></div>
    </div>
  );
}

export default App;
