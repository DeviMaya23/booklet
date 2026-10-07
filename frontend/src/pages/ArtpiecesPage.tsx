import { useState } from 'react'
import { Plus, ChevronDown } from 'lucide-react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQueryClient } from '@tanstack/react-query'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from '@/components/ui/dropdown-menu'
import { useArtpieces, ARTPIECES_QUERY_KEY } from '@/features/artpieces/api/useArtpieces'
import ArtpiecesGrid from '@/features/artpieces/components/ArtpiecesGrid'
import ArtistCharacterFilter from '@/components/ArtistCharacterFilter'
import ArtpieceFormModal from '@/features/files/components/ArtpieceFormModal'
import ArtpieceDetailView from '@/features/artpieces/components/ArtpieceDetailView'
import { useArtpiecesFilter, type SortMode } from '@/features/artpieces/hooks/useArtpiecesFilter'

const SORT_LABELS: Record<SortMode, string> = {
  newest: 'Newest first',
  alpha: 'Alphabetical',
}

export default function ArtpiecesPage() {
  const { id: selectedArtpieceId } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [modalOpen, setModalOpen] = useState(false)
  const { data: artpieces, isLoading, isError } = useArtpieces()
  const queryClient = useQueryClient()

  const {
    search, setSearch,
    sort, setSort,
    selectedArtist, setSelectedArtist,
    selectedCharacters, setSelectedCharacters,
    characterMatch, setCharacterMatch,
    filtered,
  } = useArtpiecesFilter(artpieces ?? [])

  function handleSuccess() {
    queryClient.invalidateQueries({ queryKey: ARTPIECES_QUERY_KEY })
  }

  if (selectedArtpieceId) {
    return (
      <ArtpieceDetailView
        artpieceId={selectedArtpieceId}
        onClose={() => window.history.length > 1 ? navigate(-1) : navigate('/app/artpieces')}
        onDeleted={() => navigate('/app/artpieces')}
      />
    )
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Input
          placeholder="Search artpieces..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />

        <ArtistCharacterFilter
          artist={selectedArtist}
          onArtistChange={setSelectedArtist}
          characters={selectedCharacters}
          onCharactersChange={setSelectedCharacters}
          characterMatch={characterMatch}
          onCharacterMatchChange={setCharacterMatch}
        />

        <DropdownMenu>
          <DropdownMenuTrigger className="inline-flex h-9 items-center gap-1.5 rounded-md border border-input bg-background px-3 text-sm shadow-xs hover:bg-accent hover:text-accent-foreground focus-visible:outline-none">
            {SORT_LABELS[sort]}
            <ChevronDown className="size-4" />
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onClick={() => setSort('newest')}>Newest first</DropdownMenuItem>
            <DropdownMenuItem onClick={() => setSort('alpha')}>Alphabetical</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <div className="flex-1" />

        <Button onClick={() => setModalOpen(true)}>
          <Plus className="size-4" />
          New Artpiece
        </Button>
      </div>

      {isError && (
        <p className="text-sm text-destructive">Failed to load artpieces.</p>
      )}

      {!isLoading && !isError && (
        <ArtpiecesGrid
          artpieces={filtered}
          onArtpieceOpen={(id) => navigate('/app/artpieces/' + id)}
        />
      )}

      <ArtpieceFormModal
        open={modalOpen}
        onOpenChange={setModalOpen}
        onSuccess={handleSuccess}
      />
    </div>
  )
}
