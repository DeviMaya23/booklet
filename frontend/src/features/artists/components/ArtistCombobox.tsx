import { useState } from 'react'
import { ChevronDown, X } from 'lucide-react'
import {
  Combobox,
  ComboboxInput,
  ComboboxInputGroup,
  ComboboxTrigger,
  ComboboxPopup,
  ComboboxItem,
  ComboboxEmpty,
  useComboboxFilter,
} from '@/components/ui/combobox'
import { useArtists, type Artist } from '../api/useArtists'

function isArtist(v: unknown): v is Artist {
  return !!v && typeof (v as Artist).id === 'string'
}

interface ArtistComboboxProps {
  value: Artist | null
  onChange: (artist: Artist | null) => void
  placeholder?: string
  disabled?: boolean
}

export default function ArtistCombobox({
  value,
  onChange,
  placeholder = 'Search artists…',
  disabled,
}: ArtistComboboxProps) {
  const [search, setSearch] = useState('')
  const artistsQuery = useArtists()
  const artists = artistsQuery.data ?? []
  const artistFilter = useComboboxFilter()
  const filtered = artists.filter((a) => artistFilter.contains(a, search, (a) => a.name))

  if (artistsQuery.isError) {
    return <p className="text-sm text-destructive">Failed to load artists.</p>
  }

  return (
    <Combobox
      value={value}
      onValueChange={(v) => onChange(isArtist(v) ? v : null)}
      itemToStringLabel={(a) => isArtist(a) ? a.name : ''}
      onInputValueChange={(v) => setSearch(v)}
      disabled={disabled}
    >
      <ComboboxInputGroup>
        <ComboboxInput placeholder={placeholder} />
        {value ? (
          <button
            type="button"
            className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
            onClick={() => { onChange(null); setSearch('') }}
            aria-label="Clear artist"
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
        {filtered.length === 0 && <ComboboxEmpty>No artists found</ComboboxEmpty>}
        {filtered.map((artist) => (
          <ComboboxItem key={artist.id} value={artist}>
            {artist.name}
          </ComboboxItem>
        ))}
      </ComboboxPopup>
    </Combobox>
  )
}
