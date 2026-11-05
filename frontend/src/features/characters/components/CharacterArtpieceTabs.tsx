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

function Count({ value }: { value: number }) {
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
  const { data: artpieces } = useArtpieces()
  const [viewer, setViewer] = useState<{ list: ArtpieceSummary[]; index: number } | null>(null)

  const { artworks, references } = useMemo(() => {
    const mine = (artpieces ?? []).filter((a) => a.characters.some((c) => c.id === characterId))
    return {
      artworks: mine.filter((a) => !a.is_reference),
      references: mine.filter((a) => a.is_reference),
    }
  }, [artpieces, characterId])

  function grid(list: ArtpieceSummary[]) {
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

  return (
    <TooltipProvider>
      <Tabs defaultValue="artpieces">
        <TabsList>
          <TabsTrigger value="artpieces">
            Artpieces <Count value={artworks.length} />
          </TabsTrigger>
          <TabsTrigger value="references">
            References <Count value={references.length} />
          </TabsTrigger>
        </TabsList>

        <TabsContent value="artpieces">
          {artworks.length === 0 ? (
            <EmptyState
              title={`No artpieces with ${characterName} yet`}
              hint="Add this character in an artpiece's Characters field and it will show up here."
            />
          ) : (
            grid(artworks)
          )}
        </TabsContent>

        <TabsContent value="references">
          {references.length === 0 ? (
            <EmptyState
              title={`No references with ${characterName} yet`}
              hint="Mark an artpiece as a reference and add this character to it, and it will show up here."
            />
          ) : (
            grid(references)
          )}
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
