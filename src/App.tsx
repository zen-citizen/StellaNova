import { IntlProvider } from "react-intl";
import { Sidebar, type DataSource } from "./components";
import { messages } from "./i18n";
import { BlrMap } from "./components/blr-map";

const userLocale = "en";

const dataSources: DataSource[] = [
  {
    href: "https://opencity.in/data",
    displayKey: "dataSources.openCity",
  },
  {
    href: "https://kgis.ksrsac.in/kgis/",
    displayKey: "dataSources.karnatakaGIS",
  },
  {
    href: "https://www.openstreetmap.org/about",
    displayKey: "dataSources.openStreetMap",
  },
];

const footerLinks: DataSource[] = [
  {
    href: "https://forms.gle/EmQiMpayciLdbww96",
    displayKey: "footer.error",
  },
  {
    href: "https://docs.google.com/forms/d/e/1FAIpQLScQS_-VgUFQZJedyu6iIlpoYymsKSyGUhrvPoJX1WkZGQqfLQ/viewform",
    displayKey: "footer.volunteer",
  },
  {
    href: "https://github.com/zen-citizen/StellaNova",
    displayKey: "footer.oss",
  },
];

function App() {
  return (
    <IntlProvider locale={userLocale} messages={messages[userLocale]}>
      <div className="flex justify-between">
        <div className="relative isolate h-full min-h-screen">
          {/* Sidebar  */}
          <Sidebar {...{ footerLinks, dataSources }} />
        </div>
        <div className="-z-10 h-full w-full">
          <BlrMap />
        </div>
      </div>
    </IntlProvider>
  );
}

export default App;
