import { renderHook, act, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useTagManagement } from './useTagManagement'
import { createWrapper } from '@/test/queryClient'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { error: vi.fn() } }))

import { apiFetch } from '@/lib/api'
import { toast } from 'sonner'

const sketch = { id: 'a', name: 'sketch' }
const sketches = { id: 'b', name: 'sketches' }
const tags = [sketch, sketches]

function res(status: number) {
  return { ok: status < 400, status, json: async () => ({}) } as Response
}

function setup() {
  return renderHook(() => useTagManagement(tags), { wrapper: createWrapper() })
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useTagManagement', () => {
  it('renames and leaves edit mode on success', async () => {
    vi.mocked(apiFetch).mockResolvedValue(res(200))
    const { result } = setup()

    act(() => result.current.startEdit(sketch))
    expect(result.current.editingId).toBe('a')
    await act(() => result.current.submitRename(sketch, 'doodle'))

    expect(result.current.editingId).toBeNull()
    expect(result.current.renameError).toBeNull()
  })

  it('opens the merge confirmation instead of renaming onto an existing name', async () => {
    const { result } = setup()

    act(() => result.current.startEdit(sketch))
    await act(() => result.current.submitRename(sketch, 'Sketches'))

    expect(result.current.mergeTarget).toEqual({ source: sketch, target: sketches })
    expect(apiFetch).not.toHaveBeenCalled()
    expect(result.current.editingId).toBe('a')
  })

  it('shows a conflict error when the server rejects the rename with 409', async () => {
    vi.mocked(apiFetch).mockResolvedValue(res(409))
    const { result } = setup()

    act(() => result.current.startEdit(sketch))
    await act(() => result.current.submitRename(sketch, 'doodle'))

    expect(result.current.renameError).toBe('A tag with that name already exists')
    expect(result.current.editingId).toBe('a')
  })

  it('cancelling a merge keeps the row in edit mode', async () => {
    const { result } = setup()
    act(() => result.current.startEdit(sketch))
    await act(() => result.current.submitRename(sketch, 'sketches'))

    act(() => result.current.cancelMerge())

    expect(result.current.mergeTarget).toBeNull()
    expect(result.current.editingId).toBe('a')
  })

  it('merge success closes the dialog and leaves edit mode', async () => {
    vi.mocked(apiFetch).mockResolvedValue(res(204))
    const { result } = setup()
    act(() => result.current.startEdit(sketch))
    await act(() => result.current.submitRename(sketch, 'sketches'))

    await act(() => result.current.confirmMerge())

    expect(result.current.mergeTarget).toBeNull()
    expect(result.current.editingId).toBeNull()
  })

  it('merge failure closes the dialog and toasts the error', async () => {
    vi.mocked(apiFetch).mockResolvedValue(res(404))
    const { result } = setup()
    act(() => result.current.startEdit(sketch))
    await act(() => result.current.submitRename(sketch, 'sketches'))

    await act(() => result.current.confirmMerge())

    expect(toast.error).toHaveBeenCalledWith('Tag no longer exists')
    expect(result.current.mergeTarget).toBeNull()
  })

  it('delete success closes the dialog', async () => {
    vi.mocked(apiFetch).mockResolvedValue(res(204))
    const { result } = setup()
    act(() => result.current.requestDelete(sketch))
    expect(result.current.deleteTarget).toBe(sketch)

    await act(() => result.current.confirmDelete())

    await waitFor(() => expect(result.current.deleteTarget).toBeNull())
  })

  it('delete failure keeps the dialog open and toasts', async () => {
    vi.mocked(apiFetch).mockResolvedValue(res(500))
    const { result } = setup()
    act(() => result.current.requestDelete(sketch))

    await act(() => result.current.confirmDelete())

    expect(toast.error).toHaveBeenCalledWith('Failed to delete tag')
    expect(result.current.deleteTarget).toBe(sketch)
  })
})
