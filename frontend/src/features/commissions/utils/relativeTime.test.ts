import { describe, it, expect } from 'vitest'
import { formatRelativeTime } from './relativeTime'

function daysAgo(n: number): Date {
  return new Date(Date.now() - n * 24 * 60 * 60 * 1000)
}

function daysFromNow(n: number): Date {
  return new Date(Date.now() + n * 24 * 60 * 60 * 1000)
}

describe('formatRelativeTime', () => {
  it('returns "today" for a date earlier today', () => {
    const d = new Date()
    d.setHours(d.getHours() - 1)
    expect(formatRelativeTime(d)).toBe('today')
  })

  it('returns "yesterday" for exactly 1 day ago', () => {
    expect(formatRelativeTime(daysAgo(1))).toBe('yesterday')
  })

  it('returns "N days ago" for 2–6 days ago', () => {
    expect(formatRelativeTime(daysAgo(2))).toBe('2 days ago')
    expect(formatRelativeTime(daysAgo(6))).toBe('6 days ago')
  })

  it('returns "1 week ago" for 7 days ago', () => {
    expect(formatRelativeTime(daysAgo(7))).toBe('1 week ago')
  })

  it('returns "N weeks ago" for 8–27 days ago', () => {
    expect(formatRelativeTime(daysAgo(14))).toBe('2 weeks ago')
    expect(formatRelativeTime(daysAgo(27))).toBe('3 weeks ago')
  })

  it('returns "N months ago" for 35+ days ago', () => {
    expect(formatRelativeTime(daysAgo(35))).toBe('1 month ago')
    expect(formatRelativeTime(daysAgo(90))).toBe('3 months ago')
  })

  it('returns "just now" for a future date', () => {
    expect(formatRelativeTime(daysFromNow(1))).toBe('just now')
  })
})
