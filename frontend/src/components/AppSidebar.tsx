import { NavLink, useMatch } from 'react-router-dom'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from '@/components/ui/sidebar'

const navItems = [
  { label: 'Characters', to: '/app/characters' },
  { label: 'Images', to: '/app/images' },
  { label: 'Artists', to: '/app/artists' },
]

export default function AppSidebar() {
  return (
    <Sidebar collapsible="icon">
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {navItems.map((item) => (
                <NavItem key={item.to} label={item.label} to={item.to} />
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  )
}

function NavItem({ label, to }: { label: string; to: string }) {
  const active = useMatch({ path: to, end: false })
  return (
    <SidebarMenuItem>
      <SidebarMenuButton isActive={!!active} render={<NavLink to={to} />}>
        <span>{label}</span>
      </SidebarMenuButton>
    </SidebarMenuItem>
  )
}
