import { describe, it, expect } from 'vitest'
import { findTagByName } from './findTagByName'

const tags = [
  { id: 'a', name: 'sketch' },
  { id: 'b', name: 'wip' },
]

describe('findTagByName', () => {
  it('finds an exact match', () => {
    expect(findTagByName(tags, 'sketch')).toEqual(tags[0])
  })

  it('matches regardless of case', () => {
    expect(findTagByName(tags, 'SKetch')).toEqual(tags[0])
  })

  it('ignores surrounding whitespace', () => {
    expect(findTagByName(tags, '  wip  ')).toEqual(tags[1])
  })

  it('does not match the excluded tag itself', () => {
    expect(findTagByName(tags, 'WIP', 'b')).toBeUndefined()
  })

  it('still finds a different tag when another is excluded', () => {
    expect(findTagByName(tags, 'sketch', 'b')).toEqual(tags[0])
  })

  it('returns undefined when nothing matches', () => {
    expect(findTagByName(tags, 'icon')).toBeUndefined()
  })
})
