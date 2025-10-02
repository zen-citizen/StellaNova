import SidebarContent from "./SidebarContent"
import SidebarFooter from "./SidebarFooter"
import SidebarHeader from "./SidebarHeader"

export default function Sidebar() {
  return (
    <aside className="border-border bg-surface flex h-full flex-col border-t md:border-t-0 md:border-r">
      <SidebarHeader />

      <div className="min-h-0 grow overflow-auto">
        <SidebarContent />
      </div>

      <SidebarFooter />
    </aside>
  )
}
