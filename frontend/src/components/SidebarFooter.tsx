export default function SidebarFooter() {
  return (
    <footer className="px-4 md:px-6">
      <nav className="border-t-[0.5px] border-separator py-4 flex justify-center md:justify-between gap-5">
        <a className="text-body-s text-secondary underline underline-offset-3 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer">
          Report an error
        </a>
        <a className="text-body-s text-secondary underline underline-offset-3 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer">
          Volunteer with us
        </a>
        <a className="text-body-s text-secondary underline underline-offset-3 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer">
          Open source
        </a>
      </nav>
    </footer>
  )
}
