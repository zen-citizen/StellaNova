import { clsx } from "clsx";
import { FormattedMessage } from "react-intl";
import { ExternalLink } from "./external-link/external-link";

export type DataSource = {
  href: string;
  displayName: string;
};

type SidebarProps = {
  dataSources: DataSource[];
  footerLinks: DataSource[];
};

export const Sidebar = ({ dataSources, footerLinks }: SidebarProps) => {
  return (
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
        <h1 className="text-xl text-blue-600">
          <FormattedMessage id="sidebar.title" />
        </h1>
        <ExternalLink
          href="https://zencitizen.in"
          className="text-lg text-gray-500"
        >
          <FormattedMessage id="sidebar.externalLink" />
        </ExternalLink>
      </header>
      <article className="flex flex-1 flex-col gap-y-6 border-t border-b border-gray-200 py-6 text-gray-700">
        <div className="flex flex-col gap-y-2">
          <p>
            <FormattedMessage id="sidebar.description" />
          </p>
          <p>
            <FormattedMessage
              id="sidebar.instruction"
              values={{
                b: (chunks) => <b>{chunks}</b>,
              }}
            />
          </p>
          <p className="italic">
            <FormattedMessage id="sidebar.disclaimer" />
          </p>
        </div>
        <div className="flex flex-col gap-y-2">
          <header className="text-lg font-semibold text-gray-800">
            <FormattedMessage id="sidebar.dataSources.title" />
          </header>
          <p>
            <FormattedMessage id="sidebar.dataSources.description" />
          </p>
          <ul className="mt-2 flex flex-col gap-y-3 text-sm">
            {dataSources.map(({ href, displayName }) => (
              <li key={href}>
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
                  <FormattedMessage id={displayName} />
                </ExternalLink>
              </li>
            ))}
          </ul>
        </div>
      </article>
      <footer className="py-2">
        <ul className="flex items-center justify-between text-sm text-gray-600">
          {footerLinks.map(({ href, displayName }) => (
            <ExternalLink
              key={href}
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
              <FormattedMessage id={displayName} />
            </ExternalLink>
          ))}
        </ul>
      </footer>
    </section>
  );
};
