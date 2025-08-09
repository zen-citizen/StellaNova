import { IntlProvider } from "react-intl";
import { Sidebar, type DataSource } from "./components";
import { useHtmlLang } from "./hooks/useHtmlLang";
import { messages } from "./i18n";

const dataSources: DataSource[] = [
  {
    href: "https://opencity.in/data",
    displayName: "dataSources.openCity",
  },
  {
    href: "https://kgis.ksrsac.in/kgis/",
    displayName: "dataSources.karnatakaGIS",
  },
  {
    href: "https://www.openstreetmap.org/about",
    displayName: "dataSources.openStreetMap",
  },
];

const footerLinks: DataSource[] = [
  {
    href: "https://forms.gle/EmQiMpayciLdbww96",
    displayName: "footer.error",
  },
  {
    href: "https://docs.google.com/forms/d/e/1FAIpQLScQS_-VgUFQZJedyu6iIlpoYymsKSyGUhrvPoJX1WkZGQqfLQ/viewform",
    displayName: "footer.volunteer",
  },
  {
    href: "https://github.com/zen-citizen/StellaNova",
    displayName: "footer.oss",
  },
];

function App() {
  const userLocale = useHtmlLang();
  return (
    <IntlProvider locale={userLocale} messages={messages[userLocale]}>
      <div className="relative isolate h-full min-h-screen">
        {/* Sidebar  */}
        <Sidebar {...{ footerLinks, dataSources }} />
        <div className=""></div>
      </div>
    </IntlProvider>
  );
}

export default App;
