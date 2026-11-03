import { vi } from 'vitest'
import { apiFetch } from '@/lib/api'

export const recencyCharacters = [
  { id: 'c-a', name: 'Alma' },
  { id: 'c-b', name: 'Alba' },
  { id: 'c-c', name: 'Alto' },
].map((c) => ({
  ...c,
  avatar_url: null,
  notes: null,
  is_public: false,
  verse_id: null,
  folders: [],
  created_at: '',
  updated_at: '',
}))

/**
 * Routes the mocked apiFetch: /characters* returns recencyCharacters in the
 * given (API) order, everything else returns an empty list.
 * Requires `vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))` in the test file.
 */
export function mockCharacterPickerApi() {
  vi.mocked(apiFetch).mockImplementation(async (path: string) => {
    const body = path.startsWith('/characters') ? recencyCharacters : []
    return { ok: true, json: async () => body } as Response
  })
}

export function characterRequests(): string[] {
  return vi
    .mocked(apiFetch)
    .mock.calls.map((c) => c[0] as string)
    .filter((p) => p.startsWith('/characters'))
}
