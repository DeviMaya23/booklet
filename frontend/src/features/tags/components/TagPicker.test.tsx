import { useState } from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import TagPicker from './TagPicker'
import { type Tag } from '../api/useTags'
import { createWrapper } from '@/test/queryClient'
import { apiFetch } from '@/lib/api'
import { toast } from 'sonner'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { error: vi.fn(), success: vi.fn() } }))

const outfit: Tag = { id: 't-outfit', name: 'outfit' }
const sketch: Tag = { id: 't-sketch', name: 'sketch' }
const skirt: Tag = { id: 't-skirt', name: 'skirt' }
const wip: Tag = { id: 't-wip', name: 'wip' }

type PostResult = { status: number; body?: unknown }

// Routes the mocked apiFetch: GET /tags* returns the (mutable) list in API
// order; POST /tags returns the supplied result.
function mockTagsApi(list: Tag[], post?: () => PostResult) {
  vi.mocked(apiFetch).mockImplementation(async (path: string, _getToken: unknown, init?: RequestInit) => {
    if (path.startsWith('/tags') && init?.method === 'POST') {
      const r = post?.() ?? { status: 201, body: wip }
      return { ok: r.status < 400, status: r.status, json: async () => r.body } as Response
    }
    if (path.startsWith('/tags')) return { ok: true, status: 200, json: async () => list } as Response
    return { ok: true, status: 200, json: async () => [] } as Response
  })
}

function posts() {
  return vi.mocked(apiFetch).mock.calls.filter((c) => c[2]?.method === 'POST')
}

function Harness({ initial = [], disabled, onChange }: { initial?: Tag[]; disabled?: boolean; onChange?: (t: Tag[]) => void }) {
  const [value, setValue] = useState<Tag[]>(initial)
  return (
    <TagPicker
      value={value}
      disabled={disabled}
      onChange={(t) => {
        setValue(t)
        onChange?.(t)
      }}
    />
  )
}

function renderPicker(props?: Parameters<typeof Harness>[0]) {
  const Wrapper = createWrapper()
  return render(<Harness {...props} />, { wrapper: Wrapper })
}

// The picker holds back typing until the tag list has loaded.
async function ready() {
  await waitFor(() => expect(screen.getByRole('textbox')).toBeEnabled())
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('TagPicker — chips', () => {
  it('renders selected tags as chips with remove controls', async () => {
    mockTagsApi([outfit, sketch])
    renderPicker({ initial: [outfit, sketch] })

    expect(screen.getByText('outfit')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /remove outfit/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /remove sketch/i })).toBeInTheDocument()
  })

  it('removing a chip notifies the caller with the remaining tags', async () => {
    mockTagsApi([outfit, sketch])
    const onChange = vi.fn()
    renderPicker({ initial: [outfit, sketch], onChange })

    await userEvent.click(screen.getByRole('button', { name: /remove sketch/i }))

    expect(onChange).toHaveBeenCalledWith([outfit])
  })

  it('Backspace on empty input removes the last chip', async () => {
    mockTagsApi([outfit, sketch])
    const onChange = vi.fn()
    renderPicker({ initial: [outfit, sketch], onChange })

    await userEvent.click(screen.getByRole('textbox'))
    await userEvent.keyboard('{Backspace}')

    expect(onChange).toHaveBeenCalledWith([outfit])
  })
})

describe('TagPicker — suggestions', () => {
  it('caps the focus list at 10 tags with a hint row', async () => {
    const tags: Tag[] = Array.from({ length: 25 }, (_, i) => ({ id: `t-${i + 1}`, name: `tag${String(i + 1).padStart(2, '0')}` }))
    mockTagsApi(tags)
    renderPicker()

    await ready()
    await userEvent.click(screen.getByRole('textbox'))

    const rows = await screen.findAllByRole('listitem')
    expect(rows.slice(0, 10).map((r) => r.textContent)).toEqual(tags.slice(0, 10).map((t) => t.name))
    expect(rows).toHaveLength(11)
    expect(rows[10]).toHaveTextContent(/type to search/i)
  })

  it('requests tags with sort=recent', async () => {
    mockTagsApi([outfit])
    renderPicker()

    await ready()
    expect(vi.mocked(apiFetch).mock.calls[0][0]).toBe('/tags?sort=recent')
  })

  it('lists suggestions in API order on focus with empty text', async () => {
    mockTagsApi([sketch, outfit, skirt])
    renderPicker()

    await ready()
    await userEvent.click(screen.getByRole('textbox'))

    const rows = await screen.findAllByText(/^(sketch|outfit|skirt)$/)
    expect(rows.map((r) => r.textContent)).toEqual(['sketch', 'outfit', 'skirt'])
  })

  it('filters suggestions by typed substring, case-insensitively', async () => {
    mockTagsApi([sketch, skirt, outfit])
    renderPicker()

    await ready()
    await userEvent.type(screen.getByRole('textbox'), 'SK')

    expect(await screen.findByText('sketch')).toBeInTheDocument()
    expect(screen.getByText('skirt')).toBeInTheDocument()
    expect(screen.queryByText('outfit')).not.toBeInTheDocument()
  })

  it('excludes already-selected tags from the suggestions', async () => {
    mockTagsApi([sketch, skirt])
    renderPicker({ initial: [sketch] })

    await ready()
    await userEvent.type(screen.getByRole('textbox'), 'sk')

    expect(await screen.findByText('skirt')).toBeInTheDocument()
    expect(screen.getAllByText('sketch')).toHaveLength(1) // the chip only
  })

  it('clicking a suggestion adds it and clears the text', async () => {
    mockTagsApi([sketch, skirt])
    const onChange = vi.fn()
    renderPicker({ onChange })

    await ready()
    const input = screen.getByRole('textbox')
    await userEvent.type(input, 'sk')
    await userEvent.click(await screen.findByText('sketch'))

    expect(onChange).toHaveBeenCalledWith([sketch])
    expect(input).toHaveValue('')
  })

  it('arrow keys plus Enter select a suggestion', async () => {
    mockTagsApi([sketch, skirt])
    const onChange = vi.fn()
    renderPicker({ onChange })

    await ready()
    await userEvent.type(screen.getByRole('textbox'), 'sk')
    await screen.findByText('sketch')
    await userEvent.keyboard('{ArrowDown}{Enter}')

    expect(onChange).toHaveBeenCalledWith([sketch])
  })
})

describe('TagPicker — inline creation', () => {
  it('offers a create row for an unknown name', async () => {
    mockTagsApi([sketch])
    renderPicker()

    await ready()
    await userEvent.type(screen.getByRole('textbox'), 'wip')

    expect(await screen.findByText(/Create/)).toHaveTextContent('wip')
  })

  it('does not offer a create row for an exact existing name, ignoring case', async () => {
    mockTagsApi([wip])
    renderPicker()

    await ready()
    await userEvent.type(screen.getByRole('textbox'), 'WIP')
    await screen.findByText('wip')

    expect(screen.queryByText(/Create/)).not.toBeInTheDocument()
  })

  it('choosing the create row posts the name and adds the returned tag', async () => {
    mockTagsApi([sketch], () => ({ status: 201, body: wip }))
    const onChange = vi.fn()
    renderPicker({ onChange })

    await ready()
    const input = screen.getByRole('textbox')
    await userEvent.type(input, 'wip')
    await userEvent.click(await screen.findByText(/Create/))

    await waitFor(() => expect(onChange).toHaveBeenCalledWith([wip]))
    expect(posts()).toHaveLength(1)
    expect(JSON.parse(posts()[0][2]?.body as string)).toEqual({ name: 'wip' })
    expect(input).toHaveValue('')
  })

  it('invalidates the tags list after a successful create', async () => {
    mockTagsApi([sketch], () => ({ status: 201, body: wip }))
    renderPicker()

    await ready()
    await userEvent.type(screen.getByRole('textbox'), 'wip{Enter}')

    await waitFor(() => {
      const gets = vi.mocked(apiFetch).mock.calls.filter((c) => c[2]?.method !== 'POST')
      expect(gets.length).toBeGreaterThanOrEqual(2)
    })
  })

  it('Enter on an exact existing name selects it without posting', async () => {
    mockTagsApi([wip])
    const onChange = vi.fn()
    renderPicker({ onChange })

    await ready()
    await userEvent.click(screen.getByRole('textbox'))
    await screen.findByText('wip') // list loaded
    await userEvent.type(screen.getByRole('textbox'), 'WIP{Enter}')

    await waitFor(() => expect(onChange).toHaveBeenCalledWith([wip]))
    expect(posts()).toHaveLength(0)
  })

  it('blank text creates nothing', async () => {
    mockTagsApi([sketch])
    const onChange = vi.fn()
    renderPicker({ onChange })

    await ready()
    await userEvent.type(screen.getByRole('textbox'), '   {Enter}')

    expect(posts()).toHaveLength(0)
    expect(onChange).not.toHaveBeenCalled()
  })

  it('a server error keeps the text, adds no chip and toasts', async () => {
    mockTagsApi([sketch], () => ({ status: 500 }))
    const onChange = vi.fn()
    renderPicker({ onChange })

    await ready()
    const input = screen.getByRole('textbox')
    await userEvent.type(input, 'wip{Enter}')

    await waitFor(() => expect(toast.error).toHaveBeenCalled())
    expect(onChange).not.toHaveBeenCalled()
    expect(input).toHaveValue('wip')
  })

  it('a 409 refetches and selects the existing tag without an error toast', async () => {
    const list: Tag[] = [sketch]
    mockTagsApi(list, () => {
      list.push(wip) // created elsewhere since the list was loaded
      return { status: 409 }
    })
    const onChange = vi.fn()
    renderPicker({ onChange })

    await ready()
    await userEvent.type(screen.getByRole('textbox'), 'wip{Enter}')

    await waitFor(() => expect(onChange).toHaveBeenCalledWith([wip]))
    expect(toast.error).not.toHaveBeenCalled()
  })
})

describe('TagPicker — disabled', () => {
  it('cannot be typed into', async () => {
    mockTagsApi([sketch])
    renderPicker({ disabled: true })

    const input = screen.getByRole('textbox')
    await userEvent.type(input, 'wip')

    expect(input).toBeDisabled()
    expect(input).toHaveValue('')
  })
})

describe('TagPicker — tag list loading and failure', () => {
  it('holds back input and shows a loading placeholder until the list has loaded', async () => {
    let resolve!: (r: Response) => void
    vi.mocked(apiFetch).mockImplementation(
      () => new Promise<Response>((r) => { resolve = r }),
    )
    renderPicker()

    const input = screen.getByRole('textbox')
    expect(input).toBeDisabled()
    expect(input).toHaveAttribute('placeholder', 'Loading…')

    resolve({ ok: true, status: 200, json: async () => [sketch] } as Response)
    await waitFor(() => expect(input).toBeEnabled())
    expect(input).toHaveAttribute('placeholder', 'Add tag…')
  })

  it('never offers a create row before the list has loaded', async () => {
    vi.mocked(apiFetch).mockImplementation(() => new Promise<Response>(() => {}))
    renderPicker()

    await userEvent.type(screen.getByRole('textbox'), 'wip')

    expect(screen.queryByText(/Create/)).not.toBeInTheDocument()
    expect(posts()).toHaveLength(0)
  })

  it('shows an error, disables adding, and keeps existing chips when the list fails to load', async () => {
    vi.mocked(apiFetch).mockResolvedValue({ ok: false, status: 500 } as Response)
    renderPicker({ initial: [outfit] })

    expect(await screen.findByText('Failed to load tags.')).toBeInTheDocument()
    expect(screen.getByRole('textbox')).toBeDisabled()
    expect(screen.getByText('outfit')).toBeInTheDocument()
    expect(posts()).toHaveLength(0)
  })
})
