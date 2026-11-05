import { renderHook, act, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, type Mock } from 'vitest'
import { toast } from 'sonner'
import { useCharacterDelete } from './useCharacterDelete'
import { createWrapper } from '@/test/queryClient'
import { mockApiRoutes, callsTo } from '@/test/mockApi'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

let onDeleted: Mock<() => void>

function setup() {
  return renderHook(() => useCharacterDelete({ characterId: 'c-1', onDeleted }), {
    wrapper: createWrapper(),
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  onDeleted = vi.fn()
})

describe('useCharacterDelete', () => {
  it('opens and cancels the confirmation without deleting anything', () => {
    mockApiRoutes({})
    const { result } = setup()

    act(() => result.current.askToDelete())
    expect(result.current.confirmOpen).toBe(true)
    act(() => result.current.cancel())

    expect(result.current.confirmOpen).toBe(false)
    expect(callsTo('DELETE', '/characters/c-1')).toHaveLength(0)
  })

  it('deletes on confirm, toasts and reports success', async () => {
    mockApiRoutes({})
    const { result } = setup()

    act(() => result.current.confirm())

    await waitFor(() => expect(onDeleted).toHaveBeenCalledTimes(1))
    expect(callsTo('DELETE', '/characters/c-1')).toHaveLength(1)
    expect(toast.success).toHaveBeenCalledWith('Character deleted')
  })

  it('toasts, closes the confirmation and reports nothing when the delete fails', async () => {
    mockApiRoutes({ 'DELETE /characters/c-1': { ok: false, status: 500, body: {} } })
    const { result } = setup()
    act(() => result.current.askToDelete())

    act(() => result.current.confirm())

    await waitFor(() => expect(toast.error).toHaveBeenCalledWith('Failed to delete character'))
    expect(result.current.confirmOpen).toBe(false)
    expect(onDeleted).not.toHaveBeenCalled()
  })
})
