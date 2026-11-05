import { render, screen } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { router as appRouter } from '@/router'
import AppSidebar from '@/components/AppSidebar'
import { SidebarProvider } from '@/components/ui/sidebar'
import { createWrapper } from '@/test/queryClient'
import { mockApiRoutes } from '@/test/mockApi'
import { makeCharacter, renderCharacterDetail } from '@/test/characterFixtures'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

beforeEach(() => {
  vi.clearAllMocks()
  window.matchMedia = vi.fn().mockReturnValue({
    matches: false,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
  })
  mockApiRoutes({ 'GET /characters/c-1': { body: makeCharacter() } })
})

describe('CharacterDetailPage routing', () => {
  it('renders the character for /app/characters/:id', async () => {
    renderCharacterDetail('c-1')

    expect(await screen.findByRole('heading', { name: 'Tanelle' })).toBeInTheDocument()
  })

  it('is registered in the app router at /app/characters/:id', () => {
    const paths: string[] = []
    const walk = (routes: typeof appRouter.routes) =>
      routes.forEach((r) => {
        if (r.path) paths.push(r.path)
        if (r.children) walk(r.children as typeof appRouter.routes)
      })
    walk(appRouter.routes)

    expect(paths).toContain('/app/characters/:id')
  })

  it('keeps the Characters sidebar item active on the detail route', async () => {
    const router = createMemoryRouter(
      [{ path: '/app/characters/:id', element: <AppSidebar /> }],
      { initialEntries: ['/app/characters/c-1'] },
    )
    const Wrapper = createWrapper()
    render(
      <Wrapper>
        <SidebarProvider>
          <RouterProvider router={router} />
        </SidebarProvider>
      </Wrapper>,
    )

    const characters = await screen.findByRole('link', { name: /characters/i })
    expect(characters).toHaveAttribute('data-active')
    expect(screen.getByRole('link', { name: /artpieces/i })).not.toHaveAttribute('data-active')
  })
})
