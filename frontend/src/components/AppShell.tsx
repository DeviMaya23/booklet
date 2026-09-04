import { Outlet } from 'react-router-dom'
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar'
import { TooltipProvider } from '@/components/ui/tooltip'
import AppTopBar from './AppTopBar'
import AppSidebar from './AppSidebar'

export default function AppShell() {
  return (
    <TooltipProvider>
      <SidebarProvider className="flex-col">
        <AppTopBar />
        <div className="flex min-h-0 flex-1">
          <AppSidebar />
          <SidebarInset>
            <main className="flex flex-1 flex-col p-4">
              <Outlet />
            </main>
          </SidebarInset>
        </div>
      </SidebarProvider>
    </TooltipProvider>
  )
}
