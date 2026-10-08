import { useSyncExternalStore } from 'react'

let sessionExpired = false
const listeners = new Set<() => void>()

export function subscribe(listener: () => void): () => void {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

function getSnapshot(): boolean {
  return sessionExpired
}

export function setSessionExpired(expired: boolean): void {
  if (sessionExpired === expired) return
  sessionExpired = expired
  listeners.forEach((listener) => listener())
}

export function useSessionExpired(): boolean {
  return useSyncExternalStore(subscribe, getSnapshot)
}
