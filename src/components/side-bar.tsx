import clsx from "clsx";
import { ExternalLink } from "./external-link";

export type DataSource = {
  href: string;
  displayName: string;
};

type SidebarProps = {
  dataSources: DataSource[];
  footerLinks: DataSource[];
};

export const Sidebar = ({ dataSources, footerLinks }: SidebarProps) => (
  <section
    className={clsx([
      "absolute top-0 left-0",
      "z-10",
      "flex h-full w-md flex-col",
      "bg-white shadow",
      "px-4 py-2",
    ])}
  >
    <header className="flex items-center justify-between py-2 pb-3 font-bold">
      <h1 className="text-xl text-blue-600">Civic Compass – Bengaluru</h1>
      <ExternalLink
        href="https://zencitizen.in"
        className="text-lg text-gray-500"
      >
        Zen Citizen
      </ExternalLink>
    </header>
    <article className="flex flex-1 flex-col gap-y-6 border-t border-b border-gray-200 py-6 text-gray-700">
      <div className="flex flex-col gap-y-2">
        <p>
          If you're a Bengaluru resident, you can use Civic Compass to identify
          the BBMP, BDA, Revenue, BESCOM, BWSSB offices, and Police stations for
          your area.
        </p>
        <p>
          Enter the <b>exact address</b> you need information for.
        </p>
        <p className="italic">This tool is only for Bengaluru at this time.</p>
      </div>
      <div className="flex flex-col gap-y-2">
        <header className="text-lg font-semibold text-gray-800">
          Data Sources
        </header>
        <p>
          We pull information from Government records. While we strive for
          accuracy, these sources can sometimes be incomplete or outdated.
        </p>
        <ul className="mt-2 flex flex-col gap-y-3 text-sm">
          {dataSources.map(({ href, displayName }) => (
            <li>
              <ExternalLink
                className={clsx([
                  "underline",
                  "decoration-current",
                  "hover:text-gray-900",
                  "focus-visible:text-gray-900",
                  "transition-[color,text-decoration-color]",
                  "delay-75",
                  "hover:delay-[0]",
                  "focus-visible:delay-[0]",
                ])}
                href={href}
              >
                {displayName}
              </ExternalLink>
            </li>
          ))}
        </ul>
      </div>
    </article>
    <footer className="">
      <ul className="flex items-center justify-between text-sm text-gray-600">
        {footerLinks.map(({ href, displayName }) => (
          <ExternalLink
            href={href}
            className={clsx([
              "underline",
              "decoration-current",
              "hover:text-gray-900",
              "focus-visible:text-gray-900",
              "transition-[color,text-decoration-color]",
              "delay-75",
              "hover:delay-[0]",
              "focus-visible:delay-[0]",
            ])}
          >
            {displayName}
          </ExternalLink>
        ))}
      </ul>
    </footer>
  </section>
);
