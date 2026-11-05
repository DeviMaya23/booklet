import { renderHook, act } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach, type Mock } from 'vitest'
import { toast } from 'sonner'
import { useCharacterSave } from './useCharacterSave'
import { createWrapper } from '@/test/queryClient'
import { mockApiRoutes, apiCalls, callsTo } from '@/test/mockApi'
import { makeCharacter } from '@/test/characterFixtures'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const character = makeCharacter()
const values = { name: ' Tanelle Wren ', notes: '  ', verseId: 'v-1', folderIds: ['f-a'] }
const avatarFile = new File(['jpg'], 'crop.jpg', { type: 'image/jpeg' })

let fetchMock: ReturnType<typeof vi.fn>
let reset: Mock<() => void>
let onSaved: Mock<() => void>

function setup(avatar: { localFile: File | null; avatarCleared: boolean }) {
  return renderHook(
    () =>
      useCharacterSave({
        character,
        avatar: { ...avatar, isChanged: avatar.localFile !== null || avatar.avatarCleared, reset },
        onSaved,
      }),
    { wrapper: createWrapper() },
  )
}

const routes = {
  'PUT /characters/c-1': { body: character },
  'POST /characters/c-1/avatar/init': {
    body: { id: 'up-1', upload_url: 'https://r2/put', expires_at: '' },
  },
}

beforeEach(() => {
  vi.clearAllMocks()
  reset = vi.fn()
  onSaved = vi.fn()
  fetchMock = vi.fn().mockResolvedValue({ ok: true } as Response)
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('useCharacterSave', () => {
  it('sends trimmed values, toasts and reports success when only metadata changed', async () => {
    mockApiRoutes(routes)
    const { result } = setup({ localFile: null, avatarCleared: false })

    await act(async () => {
      await result.current.save(values)
    })

    expect(callsTo('PUT', '/characters/c-1')[0].body).toEqual({
      name: 'Tanelle Wren',
      notes: null,
      is_public: false,
      folder_ids: ['f-a'],
      verse_id: 'v-1',
    })
    expect(toast.success).toHaveBeenCalledWith('Character updated')
    expect(onSaved).toHaveBeenCalledTimes(1)
    expect(reset).not.toHaveBeenCalled()
  })

  it('uploads a new avatar before the update, then resets the picker', async () => {
    mockApiRoutes(routes)
    const { result } = setup({ localFile: avatarFile, avatarCleared: false })

    await act(async () => {
      await result.current.save(values)
    })

    const order = apiCalls()
      .map((c) => `${c.method} ${c.path}`)
      .filter((p) => p.includes('/avatar') || p === 'PUT /characters/c-1')
    expect(order).toEqual([
      'POST /characters/c-1/avatar/init',
      'POST /characters/c-1/avatar/up-1/complete',
      'PUT /characters/c-1',
    ])
    expect(reset).toHaveBeenCalledTimes(1)
    expect(onSaved).toHaveBeenCalledTimes(1)
  })

  it('deletes the avatar before the update when it was removed', async () => {
    mockApiRoutes(routes)
    const { result } = setup({ localFile: null, avatarCleared: true })

    await act(async () => {
      await result.current.save(values)
    })

    const order = apiCalls()
      .map((c) => `${c.method} ${c.path}`)
      .filter((p) => p.includes('/avatar') || p === 'PUT /characters/c-1')
    expect(order).toEqual(['DELETE /characters/c-1/avatar', 'PUT /characters/c-1'])
  })

  it('stops with a toast, and does not update, when the avatar upload fails', async () => {
    mockApiRoutes(routes)
    fetchMock.mockResolvedValue({ ok: false } as Response)
    const { result } = setup({ localFile: avatarFile, avatarCleared: false })

    await act(async () => {
      await result.current.save(values)
    })

    expect(toast.error).toHaveBeenCalledWith('Failed to upload avatar')
    expect(callsTo('PUT', '/characters/c-1')).toHaveLength(0)
    expect(onSaved).not.toHaveBeenCalled()
  })

  it('stops with a toast when removing the avatar fails', async () => {
    mockApiRoutes({ ...routes, 'DELETE /characters/c-1/avatar': { ok: false, status: 500, body: {} } })
    const { result } = setup({ localFile: null, avatarCleared: true })

    await act(async () => {
      await result.current.save(values)
    })

    expect(toast.error).toHaveBeenCalledWith('Failed to remove avatar')
    expect(callsTo('PUT', '/characters/c-1')).toHaveLength(0)
  })

  it('toasts and does not report success when the update fails', async () => {
    mockApiRoutes({ 'PUT /characters/c-1': { ok: false, status: 500, body: {} } })
    const { result } = setup({ localFile: null, avatarCleared: false })

    await act(async () => {
      await result.current.save(values)
    })

    expect(toast.error).toHaveBeenCalledWith('Failed to update character')
    expect(onSaved).not.toHaveBeenCalled()
    expect(result.current.isSaving).toBe(false)
  })

  it('does nothing for a blank name', async () => {
    mockApiRoutes(routes)
    const { result } = setup({ localFile: null, avatarCleared: false })

    await act(async () => {
      await result.current.save({ ...values, name: '   ' })
    })

    expect(apiCalls()).toHaveLength(0)
    expect(onSaved).not.toHaveBeenCalled()
  })
})
