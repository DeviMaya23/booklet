import { NavLink, useMatch } from 'react-router-dom'
import { type LucideIcon, Images, UsersRound, ReceiptText, Palette, Inbox, LayoutDashboard } from 'lucide-react'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import { useFiles } from '@/features/files/api/useFiles'

const groups = [
  {
    label: 'Library',
    items: [
      { label: 'Artpieces', to: '/app/artpieces', icon: Images },
      { label: 'Characters', to: '/app/characters', icon: UsersRound },
    ],
  },
  {
    label: 'Manage',
    items: [
      { label: 'Commissions', to: '/app/commissions', icon: ReceiptText },
      { label: 'Artists', to: '/app/artists', icon: Palette },
    ],
  },
]

export default function AppSidebar() {
  const { data, isPending } = useFiles()
  const count = data?.length ?? 0

  return (
    <Sidebar collapsible="none">
      <SidebarContent>
        <div className="px-2 pt-2 flex flex-col gap-2">
          <SidebarMenu>
            <NavItem label="Dashboard" to="/app/dashboard" icon={LayoutDashboard} />
          </SidebarMenu>
          <InboxCard count={count} isPending={isPending} />
        </div>
        {groups.map((group) => (
          <SidebarGroup key={group.label}>
            <SidebarGroupLabel>{group.label}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => (
                  <NavItem key={item.to} label={item.label} to={item.to} icon={item.icon} />
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
    </Sidebar>
  )
}

function InboxCard({ count, isPending }: { count: number; isPending: boolean }) {
  const subtitle = isPending ? null : count > 0 ? `${count} to sort` : 'All sorted'

  return (
    <NavLink
      to="/app/files"
      className="flex h-14 items-center gap-3 rounded-xl bg-background px-3 shadow-xs transition-opacity hover:opacity-80"
    >
      <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-muted">
        <Inbox size={18} strokeWidth={1.75} />
      </div>
      <div className="flex min-w-0 flex-1 flex-col">
        <span className="text-[15px] font-bold leading-tight">Inbox</span>
        {subtitle !== null && (
          <span className="text-[13px] leading-tight text-muted-foreground">{subtitle}</span>
        )}
      </div>
      {!isPending && count > 0 && (
        <span
          aria-hidden="true"
          className="flex h-6 min-w-6 items-center justify-center rounded-full bg-primary px-1.5 text-[12px] font-medium text-primary-foreground"
        >
          {count}
        </span>
      )}
    </NavLink>
  )
}

function NavItem({ label, to, icon: Icon }: { label: string; to: string; icon: LucideIcon }) {
  const active = useMatch({ path: to, end: false })
  return (
    <SidebarMenuItem>
      <SidebarMenuButton
        isActive={!!active}
        className="[&[data-active]]:bg-sidebar-active-fill [&[data-active]]:font-semibold"
        render={<NavLink to={to} />}
      >
        <Icon size={18} strokeWidth={1.75} />
        <span className="text-[15px]">{label}</span>
      </SidebarMenuButton>
    </SidebarMenuItem>
  )
}
