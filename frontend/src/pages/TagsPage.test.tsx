import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import TagsPage from './TagsPage'
import { useTags } from '@/features/tags/api/useTags'
import { createWrapper } from '@/test/queryClient'
import { toast } from 'sonner'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { error: vi.fn(), success: vi.fn() } }))

import { apiFetch } from '@/lib/api'

interface FakeTag {
  id: string
  name: string
  count: number
}

// In-memory server. `recent` is the sort=recent order; the default order is by name.
let tags: FakeTag[]
let recent: string[]

function resetServer() {
  tags = [
    { id: 'outfit', name: 'outfit', count: 22 },
    { id: 'sketch', name: 'sketch', count: 14 },
    { id: 'sketches', name: 'sketches', count: 2 },
    { id: 'wip', name: 'wip', count: 0 },
  ]
  recent = ['wip', 'sketch', 'outfit', 'sketches']
}

function respond(body: unknown, status = 200) {
  return { ok: status < 400, status, json: async () => body } as Response
}

function installApi() {
  vi.mocked(apiFetch).mockImplementation(async (path: string, _t: unknown, init?: RequestInit) => {
    const method = init?.method ?? 'GET'
    if (method === 'GET' && path.startsWith('/tags')) {
      const params = new URLSearchParams(path.split('?')[1] ?? '')
      const ordered = params.get('sort') === 'recent'
        ? recent.map((id) => tags.find((t) => t.id === id)).filter((t): t is FakeTag => !!t)
        : [...tags].sort((a, b) => a.name.localeCompare(b.name))
      return respond(
        ordered.map((t) => (params.get('with_usage') ? { id: t.id, name: t.name, artpiece_count: t.count } : { id: t.id, name: t.name })),
      )
    }
    if (method === 'PATCH') {
      const id = path.split('/')[2]
      const body = JSON.parse(init?.body as string)
      tags = tags.map((t) => (t.id === id ? { ...t, name: body.name } : t))
      return respond({ id, name: body.name })
    }
    if (method === 'DELETE') {
      const id = path.split('/')[2]
      tags = tags.filter((t) => t.id !== id)
      return respond(null, 204)
    }
    if (method === 'POST' && path.endsWith('/merge')) {
      const id = path.split('/')[2]
      tags = tags.filter((t) => t.id !== id)
      return respond(null, 204)
    }
    throw new Error(`unexpected ${method} ${path}`)
  })
}

function callsWith(method: string) {
  return vi.mocked(apiFetch).mock.calls.filter((c) => (c[2]?.method ?? 'GET') === method)
}

// Stand-in for a tag picker elsewhere in the app reading the shared tag cache.
function PickerProbe() {
  const { data = [] } = useTags({ sort: 'recent' })
  return <div data-testid="picker">{data.map((t) => t.name).join(',')}</div>
}

function renderPage(withProbe = false) {
  return render(
    <MemoryRouter>
      <TagsPage />
      {withProbe && <PickerProbe />}
    </MemoryRouter>,
    { wrapper: createWrapper() },
  )
}

function rowNames() {
  return screen.getAllByRole('listitem').map((li) => within(li).getAllByText(/^[a-z]+$/)[0].textContent)
}

beforeEach(() => {
  vi.clearAllMocks()
  resetServer()
  installApi()
})

describe('TagsPage list', () => {
  it('requests usage counts and shows rows with counts, badge and links', async () => {
    renderPage()

    expect(await screen.findByText('22 artpieces →')).toBeInTheDocument()
    expect(apiFetch).toHaveBeenCalledWith('/tags?sort=recent&with_usage=true', expect.any(Function))
    expect(screen.getByText('Unused')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: '14 artpieces →' })).toHaveAttribute('href', '/app/artpieces?tags=sketch')
  })

  it('defaults to Recently used order and can switch to Name (A–Z)', async () => {
    renderPage()
    await screen.findByText('22 artpieces →')
    expect(screen.getByRole('button', { name: 'Sort tags' })).toHaveTextContent('Recently used')
    expect(rowNames()).toEqual(['wip', 'sketch', 'outfit', 'sketches'])

    await userEvent.click(screen.getByRole('button', { name: 'Sort tags' }))
    await userEvent.click(await screen.findByRole('menuitem', { name: 'Name (A–Z)' }))

    await waitFor(() => expect(rowNames()).toEqual(['outfit', 'sketch', 'sketches', 'wip']))
    expect(apiFetch).toHaveBeenCalledWith('/tags?with_usage=true', expect.any(Function))
  })

  it('combines search with sort', async () => {
    renderPage()
    await screen.findByText('22 artpieces →')

    await userEvent.type(screen.getByRole('textbox', { name: 'Search tags' }), 'SKET')
    expect(rowNames()).toEqual(['sketch', 'sketches'])

    await userEvent.click(screen.getByRole('button', { name: 'Sort tags' }))
    await userEvent.click(await screen.findByRole('menuitem', { name: 'Name (A–Z)' }))

    await waitFor(() => expect(rowNames()).toEqual(['sketch', 'sketches']))
    expect(screen.queryByText('outfit')).not.toBeInTheDocument()
  })

  it('shows neither rows nor the empty state when search matches nothing', async () => {
    renderPage()
    await screen.findByText('22 artpieces →')

    await userEvent.type(screen.getByRole('textbox', { name: 'Search tags' }), 'zzz')

    expect(screen.queryByRole('listitem')).not.toBeInTheDocument()
    expect(screen.queryByText('No tags yet')).not.toBeInTheDocument()
  })

  it('shows the empty state when the user has no tags', async () => {
    tags = []
    renderPage()

    expect(await screen.findByText('No tags yet')).toBeInTheDocument()
    expect(screen.queryByRole('listitem')).not.toBeInTheDocument()
  })
})

describe('TagsPage rename and merge', () => {
  async function startEditing(name: string) {
    renderPage()
    await screen.findByText('22 artpieces →')
    await userEvent.click(screen.getByRole('button', { name: `Edit ${name}` }))
    const input = screen.getByRole('textbox', { name: 'Tag name' })
    await userEvent.clear(input)
    return input
  }

  it('renames through PATCH when there is no collision and leaves edit mode', async () => {
    const input = await startEditing('outfit')

    await userEvent.type(input, 'outfits{Enter}')

    await waitFor(() => expect(screen.queryByRole('textbox', { name: 'Tag name' })).not.toBeInTheDocument())
    expect(callsWith('PATCH')).toHaveLength(1)
    expect(callsWith('PATCH')[0][0]).toBe('/tags/outfit')
    expect(await screen.findByText('outfits')).toBeInTheDocument()
  })

  it('treats a case-only change to its own name as a normal rename', async () => {
    const input = await startEditing('wip')

    await userEvent.type(input, 'WIP{Enter}')

    await waitFor(() => expect(callsWith('PATCH')).toHaveLength(1))
    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()
    expect(callsWith('POST')).toHaveLength(0)
  })

  it('opens the merge dialog instead of sending anything on a collision', async () => {
    const input = await startEditing('sketches')

    await userEvent.type(input, 'sketch{Enter}')

    expect(await screen.findByText('Merge into “sketch”?')).toBeInTheDocument()
    expect(callsWith('PATCH')).toHaveLength(0)
    expect(callsWith('POST')).toHaveLength(0)
  })

  it('detects a collision regardless of case', async () => {
    const input = await startEditing('sketches')

    await userEvent.type(input, 'Sketch{Enter}')

    expect(await screen.findByText('Merge into “sketch”?')).toBeInTheDocument()
    expect(callsWith('PATCH')).toHaveLength(0)
  })

  it('cancelling the merge keeps the row in edit mode and sends nothing', async () => {
    const input = await startEditing('sketches')
    await userEvent.type(input, 'sketch{Enter}')
    await screen.findByText('Merge into “sketch”?')

    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument())
    expect(screen.getByRole('textbox', { name: 'Tag name' })).toBeInTheDocument()
    expect(callsWith('POST')).toHaveLength(0)
  })

  it('confirming the merge calls the merge endpoint, refreshes the list and the shared tag cache', async () => {
    renderPage(true)
    await screen.findByText('22 artpieces →')
    expect(screen.getByTestId('picker')).toHaveTextContent('sketches')
    await userEvent.click(screen.getByRole('button', { name: 'Edit sketches' }))
    const input = screen.getByRole('textbox', { name: 'Tag name' })
    await userEvent.clear(input)
    await userEvent.type(input, 'sketch{Enter}')
    await screen.findByText('Merge into “sketch”?')

    await userEvent.click(screen.getByRole('button', { name: 'Merge tags' }))

    await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument())
    const post = callsWith('POST')[0]
    expect(post[0]).toBe('/tags/sketches/merge')
    expect(JSON.parse(post[2]?.body as string)).toEqual({ target_id: 'sketch' })
    await waitFor(() => expect(rowNames()).toEqual(['wip', 'sketch', 'outfit']))
    await waitFor(() => expect(screen.getByTestId('picker')).not.toHaveTextContent('sketches'))
  })

  it('stays in edit mode with an error when the rename request fails', async () => {
    vi.mocked(apiFetch).mockImplementation(async (_path: string, _t: unknown, init?: RequestInit) => {
      if (init?.method === 'PATCH') return respond(null, 500)
      return respond(tags.map((t) => ({ id: t.id, name: t.name, artpiece_count: t.count })))
    })
    const input = await startEditing('outfit')

    await userEvent.type(input, 'outfits{Enter}')

    expect(await screen.findByRole('alert')).toHaveTextContent('Failed to rename tag')
    expect(screen.getByRole('textbox', { name: 'Tag name' })).toBeInTheDocument()
  })
})

describe('TagsPage delete', () => {
  it('confirming deletes the tag and removes its row', async () => {
    renderPage()
    await screen.findByText('22 artpieces →')

    await userEvent.click(screen.getByRole('button', { name: 'Delete sketch' }))
    expect(await screen.findByText('Delete “sketch”?')).toBeInTheDocument()
    expect(screen.getByText(/14 artpieces use this tag/)).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Delete' }))

    await waitFor(() => expect(callsWith('DELETE')).toHaveLength(1))
    expect(callsWith('DELETE')[0][0]).toBe('/tags/sketch')
    await waitFor(() => expect(rowNames()).toEqual(['wip', 'outfit', 'sketches']))
  })

  it('cancelling closes the dialog and keeps the tag', async () => {
    renderPage()
    await screen.findByText('22 artpieces →')

    await userEvent.click(screen.getByRole('button', { name: 'Delete wip' }))
    expect(await screen.findByText(/No artpieces use this tag/)).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument())
    expect(callsWith('DELETE')).toHaveLength(0)
    expect(screen.getByText('wip')).toBeInTheDocument()
  })

  it('keeps the dialog open, toasts and refreshes the list when delete fails', async () => {
    renderPage()
    await screen.findByText('22 artpieces →')
    const base = vi.mocked(apiFetch).getMockImplementation()!
    vi.mocked(apiFetch).mockImplementation(async (path, t, init) =>
      (init?.method ?? 'GET') === 'DELETE' ? respond(null, 404) : base(path, t, init),
    )
    const listCallsBefore = callsWith('GET').length

    await userEvent.click(screen.getByRole('button', { name: 'Delete wip' }))
    await userEvent.click(await screen.findByRole('button', { name: 'Delete' }))

    await waitFor(() => expect(toast.error).toHaveBeenCalledWith('Failed to delete tag'))
    expect(screen.getByRole('alertdialog')).toBeInTheDocument()
    await waitFor(() => expect(callsWith('GET').length).toBeGreaterThan(listCallsBefore))
  })
})

describe('TagsPage load failure', () => {
  it('shows an error instead of the empty state when the list fails to load', async () => {
    vi.mocked(apiFetch).mockResolvedValue(respond(null, 500))
    renderPage()

    expect(await screen.findByText('Failed to load tags.')).toBeInTheDocument()
    expect(screen.queryByText('No tags yet')).not.toBeInTheDocument()
  })
})
