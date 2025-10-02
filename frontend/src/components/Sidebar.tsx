import SidebarContent from "./SidebarContent"
import SidebarFooter from "./SidebarFooter"
import SidebarHeader from "./SidebarHeader"

export default function Sidebar() {
  return (
    <aside className="border-separator bg-surface flex h-full flex-col">
      <SidebarHeader />

      <div className="min-h-0 grow overflow-auto">
        <SidebarContent />
      </div>

      <SidebarFooter />
    </aside>
  )
}
