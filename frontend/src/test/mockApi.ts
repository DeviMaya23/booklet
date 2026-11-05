import { vi } from 'vitest'
import { apiFetch } from '@/lib/api'

type Reply = { ok?: boolean; status?: number; body?: unknown }
type Handler = Reply | ((call: { body: unknown }) => Reply)

export interface ApiCall {
  method: string
  path: string
  body: unknown
}

/**
 * Routes the mocked `apiFetch` by "METHOD /path" (e.g. "PUT /characters/c-1").
 * Unrouted GETs return `[]`, unrouted writes return `{}`, both ok.
 * Requires `vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))` in the test file.
 */
export function mockApiRoutes(routes: Record<string, Handler>) {
  vi.mocked(apiFetch).mockImplementation(async (path: string, _getToken, init) => {
    const method = (init?.method ?? 'GET').toUpperCase()
    const handler = routes[`${method} ${path}`]
    const body = typeof init?.body === 'string' ? JSON.parse(init.body) : undefined
    const reply: Reply =
      typeof handler === 'function'
        ? handler({ body })
        : (handler ?? { body: method === 'GET' ? [] : {} })
    const ok = reply.ok ?? true
    return {
      ok,
      status: reply.status ?? (ok ? 200 : 500),
      json: async () => reply.body,
    } as Response
  })
}

export function apiCalls(): ApiCall[] {
  return vi.mocked(apiFetch).mock.calls.map(([path, , init]) => ({
    method: (init?.method ?? 'GET').toUpperCase(),
    path,
    body: typeof init?.body === 'string' ? JSON.parse(init.body) : undefined,
  }))
}

export function callsTo(method: string, path: string): ApiCall[] {
  return apiCalls().filter((c) => c.method === method && c.path === path)
}
