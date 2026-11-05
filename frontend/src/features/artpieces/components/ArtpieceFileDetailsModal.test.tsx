import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect, beforeEach } from 'vitest'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import ArtpieceFileDetailsModal from './ArtpieceFileDetailsModal'
import { type ArtpieceFile } from '../api/useArtpiece'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: () => ({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))

import { apiFetch } from '@/lib/api'

function makeFile(id: string, overrides: Partial<ArtpieceFile> = {}): ArtpieceFile {
  return {
    id,
    file_url: `https://cdn/${id}`,
    thumbnail_url: null,
    mime_type: 'image/png',
    name: `name-${id}`,
    notes: `notes-${id}`,
    is_visible: true,
    extension: 'png',
    width: null,
    height: null,
    ...overrides,
  }
}

const files = [makeFile('a'), makeFile('b', { is_visible: false, notes: null })]

let qc: QueryClient

function renderModal(f: ArtpieceFile[] = files) {
  qc = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
  const invalidate = vi.spyOn(qc, 'invalidateQueries')
  render(
    <QueryClientProvider client={qc}>
      <ArtpieceFileDetailsModal open onOpenChange={() => {}} artpieceId="ap1" files={f} />
    </QueryClientProvider>,
  )
  return { invalidate }
}

function lastBody(): Record<string, unknown> {
  const calls = vi.mocked(apiFetch).mock.calls
  return JSON.parse((calls[calls.length - 1][2] as RequestInit).body as string)
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(apiFetch).mockResolvedValue({ ok: true } as Response)
})

describe('ArtpieceFileDetailsModal', () => {
  it('renders a row for every attached file, including hidden ones, pre-filled', () => {
    renderModal()
    expect(screen.getByLabelText('Name for name-a')).toHaveValue('name-a')
    expect(screen.getByLabelText('Notes for name-a')).toHaveValue('notes-a')
    expect(screen.getByLabelText('Visible: name-a')).toBeChecked()
    expect(screen.getByLabelText('Name for name-b')).toHaveValue('name-b')
    expect(screen.getByLabelText('Notes for name-b')).toHaveValue('')
    expect(screen.getByLabelText('Visible: name-b')).not.toBeChecked()
    expect(screen.getByText('Hidden')).toBeInTheDocument()
  })

  it('shows an empty state when no files are attached', () => {
    renderModal([])
    expect(screen.getByText('No files attached.')).toBeInTheDocument()
  })

  it.each([
    ['name', 'Name for name-a', 'Notes for name-a'],
    ['notes', 'Notes for name-a', 'Name for name-a'],
  ])('shows an inline error when saving %s fails, then clears it on success', async (_field, label, otherLabel) => {
    const user = userEvent.setup()
    renderModal()
    vi.mocked(apiFetch).mockResolvedValueOnce({ ok: false } as Response)

    const input = screen.getByLabelText(label)
    await user.type(input, 'x')
    await user.tab()
    expect(await screen.findByText("Couldn't save. Try again.")).toBeInTheDocument()
    // the other field is unaffected
    expect(screen.getAllByText("Couldn't save. Try again.")).toHaveLength(1)
    expect(screen.getByLabelText(otherLabel)).not.toHaveClass('border-destructive')

    await user.click(input)
    await user.tab()
    await waitFor(() => expect(screen.queryByText("Couldn't save. Try again.")).not.toBeInTheDocument())
  })

  it('does not send a failed notes draft when a different field is saved later', async () => {
    const user = userEvent.setup()
    renderModal()
    vi.mocked(apiFetch).mockResolvedValueOnce({ ok: false } as Response)

    await user.type(screen.getByLabelText('Notes for name-a'), '!')
    await user.tab()
    await screen.findByText("Couldn't save. Try again.")

    await user.type(screen.getByLabelText('Name for name-a'), '?')
    await user.tab()

    await waitFor(() => expect(apiFetch).toHaveBeenCalledTimes(2))
    expect(lastBody()).toEqual({ name: 'name-a?', notes: 'notes-a' })
  })

  it('does not call the API when a field is blurred unchanged', async () => {
    const user = userEvent.setup()
    renderModal()

    await user.click(screen.getByLabelText('Name for name-a'))
    await user.tab()

    expect(apiFetch).not.toHaveBeenCalled()
  })

  it('toggles visibility immediately and invalidates the artpiece query', async () => {
    const user = userEvent.setup()
    const { invalidate } = renderModal()

    await user.click(screen.getByLabelText('Visible: name-a'))

    await waitFor(() => expect(screen.getByLabelText('Visible: name-a')).not.toBeChecked())
    expect(lastBody()).toMatchObject({ is_visible: false })
    await waitFor(() =>
      expect(invalidate).toHaveBeenCalledWith({ queryKey: ['artpieces', 'ap1'] }),
    )
  })

  it('still invalidates the artpiece query when the modal closes before the save returns', async () => {
    const user = userEvent.setup()
    let resolveSave: (r: Response) => void = () => {}
    vi.mocked(apiFetch).mockReturnValueOnce(new Promise<Response>((r) => { resolveSave = r }))
    const qcLocal = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })
    const invalidate = vi.spyOn(qcLocal, 'invalidateQueries')
    const { rerender } = render(
      <QueryClientProvider client={qcLocal}>
        <ArtpieceFileDetailsModal open onOpenChange={() => {}} artpieceId="ap1" files={files} />
      </QueryClientProvider>,
    )

    await user.click(screen.getByLabelText('Visible: name-a'))
    rerender(
      <QueryClientProvider client={qcLocal}>
        <ArtpieceFileDetailsModal open={false} onOpenChange={() => {}} artpieceId="ap1" files={files} />
      </QueryClientProvider>,
    )
    resolveSave({ ok: true } as Response)

    await waitFor(() =>
      expect(invalidate).toHaveBeenCalledWith({ queryKey: ['artpieces', 'ap1'] }),
    )
  })

  it('reverts the toggle and shows an error when the visibility save fails', async () => {
    const user = userEvent.setup()
    const { invalidate } = renderModal()
    vi.mocked(apiFetch).mockResolvedValueOnce({ ok: false } as Response)

    await user.click(screen.getByLabelText('Visible: name-a'))

    expect(await screen.findByText("Couldn't save. Try again.")).toBeInTheDocument()
    expect(screen.getByLabelText('Visible: name-a')).toBeChecked()
    expect(invalidate).not.toHaveBeenCalled()
  })
})
