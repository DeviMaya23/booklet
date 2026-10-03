import { useState } from 'react'
import { ChevronDown, Plus } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useFiles } from '@/features/files/api/useFiles'
import FileInboxGrid from '@/features/files/components/FileInboxGrid'
import ArtpieceFormModal from '@/features/files/components/ArtpieceFormModal'

export default function FilesPage() {
  const [search, setSearch] = useState('')
  const [selection, setSelection] = useState<Set<string>>(new Set())
  const [artpieceModalOpen, setArtpieceModalOpen] = useState(false)
  const [artpieceModalFileIds, setArtpieceModalFileIds] = useState<string[]>([])
  const { data: files = [], isLoading } = useFiles()

  function openNewArtpieceModal(fileIds: string[]) {
    setArtpieceModalFileIds(fileIds)
    setArtpieceModalOpen(true)
  }

  const filtered = search.trim()
    ? files.filter((f) =>
        f.notes?.toLowerCase().includes(search.toLowerCase()) ||
        f.name?.toLowerCase().includes(search.toLowerCase())
      )
    : files

  return (
    <div className="flex flex-col gap-4 h-full relative">
      <div className="flex items-center gap-2">
        <Input
          placeholder="Search files..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
        {selection.size > 0 && (
          <DropdownMenu>
            <DropdownMenuTrigger className="ml-auto inline-flex items-center gap-1 rounded-md border px-3 py-2 text-sm font-medium hover:bg-accent">
              Add to artpiece <ChevronDown className="h-4 w-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => openNewArtpieceModal(Array.from(selection))}>
                New Artpiece
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => {/* no-op */}}>
                Add to Existing Artpiece
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </div>

      {!isLoading && (
        <FileInboxGrid
          files={filtered}
          selection={selection}
          onSelectionChange={setSelection}
          onNewArtpiece={(fileIds) => openNewArtpieceModal(fileIds)}
        />
      )}

      <p className="text-center text-sm text-muted-foreground">drag files here to upload</p>

      <Button
        size="icon"
        className="fixed bottom-6 right-6 h-12 w-12 rounded-full shadow-lg"
        onClick={() => {/* no-op */}}
        aria-label="Upload file"
      >
        <Plus className="h-5 w-5" />
      </Button>

      <ArtpieceFormModal
        open={artpieceModalOpen}
        onOpenChange={setArtpieceModalOpen}
        initialFileIds={artpieceModalFileIds}
        onSuccess={() => setSelection(new Set())}
      />
    </div>
  )
}
