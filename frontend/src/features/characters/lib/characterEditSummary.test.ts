import { describe, it, expect } from 'vitest'
import { characterEditSummary, type CharacterEditSummaryInput } from './characterEditSummary'

const character = {
  id: 'c-1',
  name: 'Tanelle',
  avatar_url: null,
  notes: 'Apothecary from Dyffrynne',
  is_public: false,
  verse_id: 'v-1',
  folders: [
    { id: 'f-a', name: 'Plant' },
    { id: 'f-b', name: 'Cottagecore' },
    { id: 'f-c', name: 'Old refs' },
  ],
  created_at: '',
  updated_at: '',
}

function input(overrides: Partial<CharacterEditSummaryInput> = {}): CharacterEditSummaryInput {
  return {
    character,
    name: 'Tanelle',
    notes: 'Apothecary from Dyffrynne',
    verseId: 'v-1',
    avatar: { hasNewFile: false, cleared: false },
    selectedFolderIds: ['f-a', 'f-b', 'f-c'],
    ...overrides,
  }
}

describe('characterEditSummary', () => {
  it('is null when nothing differs', () => {
    expect(characterEditSummary(input())).toBeNull()
  })

  it('lists name, verse and notes changes', () => {
    expect(
      characterEditSummary(input({ name: 'Tan', verseId: 'v-2', notes: 'New notes' })),
    ).toBe('Unsaved changes: name changed, verse changed, notes changed')
  })

  it('ignores surrounding whitespace in name and notes', () => {
    expect(characterEditSummary(input({ name: '  Tanelle ', notes: ' Apothecary from Dyffrynne  ' }))).toBeNull()
  })

  it('treats empty notes as unchanged when the character has none', () => {
    const noNotes = { ...character, notes: null }
    expect(characterEditSummary(input({ character: noNotes, notes: '' }))).toBeNull()
  })

  it('drops "verse changed" when the original verse is selected again', () => {
    const changed = characterEditSummary(input({ verseId: 'v-2' }))
    expect(changed).toBe('Unsaved changes: verse changed')

    expect(characterEditSummary(input({ verseId: 'v-1' }))).toBeNull()
  })

  it('reports clearing the verse', () => {
    expect(characterEditSummary(input({ verseId: null }))).toBe('Unsaved changes: verse changed')
  })

  it('reports a new avatar and a removed avatar', () => {
    expect(characterEditSummary(input({ avatar: { hasNewFile: true, cleared: false } }))).toBe(
      'Unsaved changes: avatar changed',
    )
    expect(characterEditSummary(input({ avatar: { hasNewFile: false, cleared: true } }))).toBe(
      'Unsaved changes: avatar removed',
    )
  })

  it('counts added and removed folders with singular and plural wording', () => {
    expect(characterEditSummary(input({ selectedFolderIds: ['f-a', 'f-b'] }))).toBe(
      'Unsaved changes: 1 folder removed',
    )
    expect(characterEditSummary(input({ selectedFolderIds: ['f-a'] }))).toBe(
      'Unsaved changes: 2 folders removed',
    )
    expect(
      characterEditSummary(input({ selectedFolderIds: ['f-a', 'f-b', 'f-c', 'f-d', 'f-e'] })),
    ).toBe('Unsaved changes: 2 folders added')
  })

  it('does not count folders Bookleaf dropped as a change', () => {
    expect(
      characterEditSummary(input({ selectedFolderIds: ['f-a', 'f-b'], staleFolderIds: ['f-c'] })),
    ).toBeNull()
  })

  it('still counts removals the user made alongside stale folders', () => {
    expect(
      characterEditSummary(input({ selectedFolderIds: ['f-a'], staleFolderIds: ['f-c'] })),
    ).toBe('Unsaved changes: 1 folder removed')
  })

  it('combines changes in the documented order', () => {
    expect(
      characterEditSummary(input({ verseId: 'v-2', selectedFolderIds: ['f-a', 'f-b', 'f-new'] })),
    ).toBe('Unsaved changes: verse changed, 1 folder added, 1 folder removed')
  })
})
