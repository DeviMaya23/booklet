import { useState } from 'react'
import { Plus } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useArtists, type Artist } from '@/features/artists/api/useArtists'
import ArtistsList from '@/features/artists/components/ArtistsList'
import ArtistFormModal from '@/features/artists/components/ArtistFormModal'

export default function ArtistsPage() {
  const [search, setSearch] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [selectedArtist, setSelectedArtist] = useState<Artist | undefined>(undefined)

  const { data: artists = [], isLoading } = useArtists()

  const filtered = search.trim()
    ? artists.filter((artist) =>
        artist.name.toLowerCase().includes(search.toLowerCase())
      )
    : artists

  function handleNewClick() {
    setSelectedArtist(undefined)
    setModalOpen(true)
  }

  function handleEditClick(artist: Artist) {
    setSelectedArtist(artist)
    setModalOpen(true)
  }

  function handleModalOpenChange(open: boolean) {
    setModalOpen(open)
    if (!open) setSelectedArtist(undefined)
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Input
          placeholder="Search artists..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
        <Button variant="outline" className="ml-auto" onClick={handleNewClick}>
          <Plus />
          New
        </Button>
      </div>

      {isLoading ? null : (
        <ArtistsList artists={filtered} onEditClick={handleEditClick} />
      )}

      <ArtistFormModal
        key={selectedArtist?.id ?? 'new'}
        open={modalOpen}
        onOpenChange={handleModalOpenChange}
        artist={selectedArtist}
      />
    </div>
  )
}
