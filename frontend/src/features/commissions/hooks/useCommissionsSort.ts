import { useState } from 'react'
import { type Commission } from '../api/useCommissions'
import { timeTakenDays } from '../utils/timeTakenDays'

export type SortKey = 'title' | 'artist_name' | 'status' | 'paid' | 'paid_date' | 'last_contacted_at' | 'time_taken'
export type SortDir = 'asc' | 'desc'

export function useCommissionsSort(commissions: Commission[]) {
  const [sortKey, setSortKey] = useState<SortKey>('title')
  const [sortDir, setSortDir] = useState<SortDir>('asc')

  function toggleSort(key: SortKey) {
    if (sortKey === key) {
      setSortDir(d => (d === 'asc' ? 'desc' : 'asc'))
    } else {
      setSortKey(key)
      setSortDir('asc')
    }
  }

  const sorted = [...commissions].sort((a, b) => {
    let aVal: string | number | null = null
    let bVal: string | number | null = null

    if (sortKey === 'title') { aVal = a.title ?? ''; bVal = b.title ?? '' }
    else if (sortKey === 'artist_name') { aVal = a.artist_name ?? ''; bVal = b.artist_name ?? '' }
    else if (sortKey === 'status') { aVal = a.status; bVal = b.status }
    else if (sortKey === 'paid') { aVal = a.paid ? 1 : 0; bVal = b.paid ? 1 : 0 }
    else if (sortKey === 'paid_date') { aVal = a.paid_date ?? ''; bVal = b.paid_date ?? '' }
    else if (sortKey === 'last_contacted_at') { aVal = a.last_contacted_at ?? ''; bVal = b.last_contacted_at ?? '' }
    else if (sortKey === 'time_taken') { aVal = timeTakenDays(a) ?? -1; bVal = timeTakenDays(b) ?? -1 }

    if (aVal == null || bVal == null) return 0
    if (aVal < bVal) return sortDir === 'asc' ? -1 : 1
    if (aVal > bVal) return sortDir === 'asc' ? 1 : -1
    return 0
  })

  return { sorted, sortKey, sortDir, toggleSort }
}
