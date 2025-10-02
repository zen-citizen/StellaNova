export default function SidebarHeader() {
  return (
    <header className="px-4 md:px-6">
      <div className="border-separator border-b-[0.5px] py-4 flex flex-col items-center gap-2 md:flex-row md:items-start md:justify-between">
        <h1 className="text-heading-xs md:text-heading-s text-accent font-bold">
          Civic Compass — Bengaluru
        </h1>
        <a
          className="text-body-m md:text-body-l text-secondary font-semibold underline underline-offset-2 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer"
          href="https://zencitizen.in/"
          target="_blank"
        >
          Zen Citizen
        </a>
      </div>
    </header>
  )
}
