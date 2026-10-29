import { useState } from 'react'
import { ChevronDown, X } from 'lucide-react'
import {
  Combobox,
  ComboboxList,
  ComboboxInput,
  ComboboxInputGroup,
  ComboboxTrigger,
  ComboboxPopup,
  ComboboxItem,
  ComboboxEmpty,
  useComboboxFilter,
} from '@/components/ui/combobox'
import { useVerses, type Verse } from '../api/useVerses'

function isVerse(v: unknown): v is Verse {
  return !!v && typeof (v as Verse).id === 'string'
}

function sortVerses(verses: Verse[]): Verse[] {
  return [...verses].sort((a, b) => {
    if (a.last_used_at && b.last_used_at) {
      return new Date(b.last_used_at).getTime() - new Date(a.last_used_at).getTime()
    }
    if (a.last_used_at) return -1
    if (b.last_used_at) return 1
    return a.name.localeCompare(b.name)
  })
}

interface VerseComboboxProps {
  value: Verse | null
  onChange: (verse: Verse | null) => void
  placeholder?: string
  disabled?: boolean
  inputId?: string
}

export default function VerseCombobox({
  value,
  onChange,
  placeholder = 'Search verses…',
  disabled,
  inputId,
}: VerseComboboxProps) {
  const [search, setSearch] = useState('')
  const versesQuery = useVerses()
  const verses = versesQuery.data ?? []
  const verseFilter = useComboboxFilter()
  const filtered = sortVerses(verses).filter((v) => verseFilter.contains(v, search, (v) => v.name))

  if (versesQuery.isError) {
    return <p className="text-sm text-destructive">Failed to load verses.</p>
  }

  return (
    <Combobox
      value={value}
      onValueChange={(v) => onChange(isVerse(v) ? v : null)}
      itemToStringLabel={(v) => (isVerse(v) ? v.name : '')}
      onInputValueChange={(v) => setSearch(v)}
      disabled={disabled}
    >
      <ComboboxInputGroup>
        <ComboboxInput id={inputId} placeholder={placeholder} />
        {value ? (
          <button
            type="button"
            className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
            onClick={() => {
              onChange(null)
              setSearch('')
            }}
            aria-label="Clear verse"
          >
            <X className="size-4" />
          </button>
        ) : (
          <ComboboxTrigger>
            <ChevronDown className="size-4" />
          </ComboboxTrigger>
        )}
      </ComboboxInputGroup>
      <ComboboxPopup>
        <ComboboxList>
          {filtered.length === 0 && <ComboboxEmpty>No verses found</ComboboxEmpty>}
          {filtered.map((verse) => (
            <ComboboxItem key={verse.id} value={verse}>
              {verse.name}
            </ComboboxItem>
          ))}
        </ComboboxList>
      </ComboboxPopup>
    </Combobox>
  )
}
