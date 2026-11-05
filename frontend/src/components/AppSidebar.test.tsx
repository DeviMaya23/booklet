import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, it, expect, vi, beforeAll } from 'vitest'
import { SidebarProvider } from '@/components/ui/sidebar'
import AppSidebar from './AppSidebar'

vi.mock('@/features/files/api/useFiles', () => ({
  useFiles: () => ({ data: [], isPending: false }),
}))

beforeAll(() => {
  // jsdom has no matchMedia; SidebarProvider's mobile hook needs it.
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
  }))
})

function renderSidebar(path = '/app/dashboard') {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <SidebarProvider>
        <AppSidebar />
        <Routes>
          <Route path="/app/tags" element={<div>Tags route</div>} />
          <Route path="*" element={<div>Other route</div>} />
        </Routes>
      </SidebarProvider>
    </MemoryRouter>,
  )
}

describe('AppSidebar', () => {
  it('lists Commissions, Artists and Tags in order under Manage', () => {
    renderSidebar()

    const label = screen.getByText('Manage')
    const group = label.parentElement as HTMLElement
    const links = within(group).getAllByRole('link').map((l) => l.textContent)
    expect(links).toEqual(['Commissions', 'Artists', 'Tags'])
  })

  it('marks Tags active on /app/tags', () => {
    renderSidebar('/app/tags')

    const tags = screen.getByRole('link', { name: 'Tags' })
    expect(tags).toHaveAttribute('data-active')
    expect(screen.getByRole('link', { name: 'Artists' })).not.toHaveAttribute('data-active')
  })

  it('navigates to /app/tags from the Tags item', async () => {
    renderSidebar()

    await userEvent.click(screen.getByRole('link', { name: 'Tags' }))

    expect(screen.getByText('Tags route')).toBeInTheDocument()
  })
})
