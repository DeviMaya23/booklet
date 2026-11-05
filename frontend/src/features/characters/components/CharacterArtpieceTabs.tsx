import { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { useArtpieces, type ArtpieceSummary } from '@/features/artpieces/api/useArtpieces'
import ArtpiecesGrid from '@/features/artpieces/components/ArtpiecesGrid'
import { toViewerFiles } from '@/features/artpieces/lib/toViewerFiles'
import FileViewer from '@/features/files/components/FileViewer'

interface CharacterArtpieceTabsProps {
  characterId: string
  characterName: string
}

function Count({ value }: { value: number | null }) {
  if (value === null) return null
  return (
    <span className="rounded-md bg-background/60 px-1.5 text-xs font-normal text-muted-foreground">
      {value}
    </span>
  )
}

function EmptyState({ title, hint }: { title: string; hint: string }) {
  return (
    <div className="flex flex-col items-center gap-1 rounded-xl border border-dashed px-6 py-12 text-center">
      <p className="font-semibold">{title}</p>
      <p className="text-sm text-muted-foreground">{hint}</p>
    </div>
  )
}

export default function CharacterArtpieceTabs({ characterId, characterName }: CharacterArtpieceTabsProps) {
  const navigate = useNavigate()
  const { data: artpieces, isPending } = useArtpieces()
  const [viewer, setViewer] = useState<{ list: ArtpieceSummary[]; index: number } | null>(null)

  const { artworks, references } = useMemo(() => {
    const mine = (artpieces ?? []).filter((a) => a.characters.some((c) => c.id === characterId))
    return {
      artworks: mine.filter((a) => !a.is_reference),
      references: mine.filter((a) => a.is_reference),
    }
  }, [artpieces, characterId])

  // No data yet: still loading, or the request failed. Never show that as "empty".
  const unavailable = artpieces === undefined

  function panel(list: ArtpieceSummary[], empty: { title: string; hint: string }) {
    if (unavailable) {
      return isPending ? (
        <p className="text-sm text-muted-foreground">Loading…</p>
      ) : (
        <p className="text-sm text-destructive">Failed to load artpieces.</p>
      )
    }
    if (list.length === 0) return <EmptyState {...empty} />
    return (
      <ArtpiecesGrid
        artpieces={list}
        tileSize="medium"
        showDetails
        onTileClick={(_id, index) => setViewer({ list, index })}
        onInfoClick={(id) => navigate(`/app/artpieces/${id}`)}
      />
    )
  }

  const artworkCount = unavailable ? null : artworks.length
  const referenceCount = unavailable ? null : references.length

  return (
    <TooltipProvider>
      <Tabs defaultValue="artpieces">
        <TabsList>
          <TabsTrigger value="artpieces">
            Artpieces <Count value={artworkCount} />
          </TabsTrigger>
          <TabsTrigger value="references">
            References <Count value={referenceCount} />
          </TabsTrigger>
        </TabsList>

        <TabsContent value="artpieces">
          {panel(artworks, {
            title: `No artpieces with ${characterName} yet`,
            hint: "Add this character in an artpiece's Characters field and it will show up here.",
          })}
        </TabsContent>

        <TabsContent value="references">
          {panel(references, {
            title: `No references with ${characterName} yet`,
            hint: 'Mark an artpiece as a reference and add this character to it, and it will show up here.',
          })}
        </TabsContent>
      </Tabs>

      {viewer && (
        <FileViewer
          files={toViewerFiles(viewer.list)}
          initialIndex={viewer.index}
          onClose={() => setViewer(null)}
        />
      )}
    </TooltipProvider>
  )
}
