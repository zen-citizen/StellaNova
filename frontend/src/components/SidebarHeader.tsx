export default function SidebarHeader() {
  return (
    <header className="px-4 md:px-6">
      <div className="border-separator border-b-[0.5px] py-4 flex flex-col items-center gap-2 md:flex-row md:items-baseline md:justify-between">
        <h1 className="text-heading-xs md:text-heading-s text-accent font-bold">
          Civic Compass — Bengaluru
        </h1>
        <div className="text-body-m md:text-body-l text-secondary font-semibold">
          Zen Citizen
        </div>
      </div>
    </header>
  )
}
