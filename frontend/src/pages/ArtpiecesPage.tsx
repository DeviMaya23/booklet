import { useState } from 'react'
import { Plus, ChevronDown, SlidersHorizontal } from 'lucide-react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQueryClient } from '@tanstack/react-query'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from '@/components/ui/dropdown-menu'
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from '@/components/ui/popover'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Skeleton } from '@/components/ui/skeleton'
import { useArtpieces, ARTPIECES_QUERY_KEY } from '@/features/artpieces/api/useArtpieces'
import ArtpiecesGrid, { type TileSize } from '@/features/artpieces/components/ArtpiecesGrid'
import ArtistCharacterFilter from '@/components/ArtistCharacterFilter'
import ArtpieceFormModal from '@/features/files/components/ArtpieceFormModal'
import ArtpieceDetailView from '@/features/artpieces/components/ArtpieceDetailView'
import { useArtpiecesFilter, type SortMode } from '@/features/artpieces/hooks/useArtpiecesFilter'
import FileViewer, { type ViewerFile } from '@/features/files/components/FileViewer'

const SORT_LABELS: Record<SortMode, string> = {
  newest: 'Newest first',
  alpha: 'Alphabetical',
}

const TILE_SIZES: TileSize[] = ['small', 'medium', 'large']
const TILE_SIZE_LABELS: Record<TileSize, string> = {
  small: 'Small',
  medium: 'Medium',
  large: 'Large',
}

function ArtpiecesGridSkeleton({ tileSize }: { tileSize: TileSize }) {
  const cols = {
    small: 'grid-cols-4 sm:grid-cols-6 md:grid-cols-8',
    medium: 'grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5',
    large: 'grid-cols-1 sm:grid-cols-2 md:grid-cols-3',
  }[tileSize]
  return (
    <div className={`grid gap-3 ${cols}`}>
      {Array.from({ length: 12 }).map((_, i) => (
        <Skeleton key={i} className="aspect-square rounded-lg" />
      ))}
    </div>
  )
}

function readStoredTileSize(): TileSize {
  try {
    const v = localStorage.getItem('gallery.tileSize')
    if (v === 'small' || v === 'medium' || v === 'large') return v
  } catch { /* storage blocked in private mode */ }
  return 'medium'
}

function readStoredShowDetails(): boolean {
  try {
    const v = localStorage.getItem('gallery.showDetails')
    if (v === 'false') return false
  } catch { /* storage blocked in private mode */ }
  return true
}

export default function ArtpiecesPage() {
  const { id: selectedArtpieceId } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [modalOpen, setModalOpen] = useState(false)
  const [tileSize, setTileSize] = useState<TileSize>(readStoredTileSize)
  const [showDetails, setShowDetails] = useState<boolean>(readStoredShowDetails)
  const [viewerIndex, setViewerIndex] = useState<number | null>(null)

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

  function handleTileSizeChange(size: TileSize) {
    setTileSize(size)
    try { localStorage.setItem('gallery.tileSize', size) } catch { /* storage blocked */ }
  }

  function handleShowDetailsChange(checked: boolean) {
    setShowDetails(checked)
    try { localStorage.setItem('gallery.showDetails', String(checked)) } catch { /* storage blocked */ }
  }

  const viewerFiles: ViewerFile[] = filtered.map((a) => ({
    id: a.cover_file_id ?? a.id,
    name:  a.title ?? a.cover_file_name ?? 'Untitled',
    mimeType: a.cover_file_mime_type ?? 'application/octet-stream',
    thumbnailUrl: a.thumbnail_url,
    previewUrl: a.cover_file_url ?? '',
  }))

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
    <TooltipProvider>
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

          <Popover>
            <PopoverTrigger
              className="inline-flex h-9 items-center gap-1.5 rounded-md border border-input bg-background px-3 text-sm shadow-xs hover:bg-accent hover:text-accent-foreground focus-visible:outline-none"
              aria-label="View options"
            >
              <SlidersHorizontal className="size-4" />
              View
            </PopoverTrigger>
            <PopoverContent align="end" className="w-52">
              <div className="flex flex-col gap-3">
                <div>
                  <p className="mb-1.5 text-xs font-medium text-muted-foreground">Tile size</p>
                  <div className="flex overflow-hidden rounded-md border border-input">
                    {TILE_SIZES.map((size) => (
                      <button
                        key={size}
                        type="button"
                        onClick={() => handleTileSizeChange(size)}
                        className={`flex-1 py-1.5 text-xs transition-colors ${
                          tileSize === size
                            ? 'bg-primary text-primary-foreground'
                            : 'hover:bg-accent hover:text-accent-foreground'
                        }`}
                      >
                        {TILE_SIZE_LABELS[size]}
                      </button>
                    ))}
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <label htmlFor="show-details-switch" className="text-sm">
                    Show details
                  </label>
                  <Switch
                    id="show-details-switch"
                    checked={showDetails}
                    onCheckedChange={handleShowDetailsChange}
                  />
                </div>
              </div>
            </PopoverContent>
          </Popover>

          <div className="flex-1" />

          <Button onClick={() => setModalOpen(true)}>
            <Plus className="size-4" />
            New Artpiece
          </Button>
        </div>

        {isError && (
          <p className="text-sm text-destructive">Failed to load artpieces.</p>
        )}

        {isLoading && (
          <ArtpiecesGridSkeleton tileSize={tileSize} />
        )}

        {!isLoading && !isError && (
          <ArtpiecesGrid
            artpieces={filtered}
            tileSize={tileSize}
            showDetails={showDetails}
            onTileClick={(_id, index) => setViewerIndex(index)}
            onInfoClick={(id) => navigate(`/app/artpieces/${id}`)}
          />
        )}

        {viewerIndex !== null && (
          <FileViewer
            files={viewerFiles}
            initialIndex={viewerIndex}
            onClose={() => setViewerIndex(null)}
          />
        )}

        <ArtpieceFormModal
          open={modalOpen}
          onOpenChange={setModalOpen}
          onSuccess={handleSuccess}
        />
      </div>
    </TooltipProvider>
  )
}
