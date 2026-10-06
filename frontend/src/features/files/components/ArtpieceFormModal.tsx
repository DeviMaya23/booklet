import { useCallback, useEffect, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { Plus } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import TokenInput from '@/components/TokenInput'
import { type Artist } from '@/features/artists/api/useArtists'
import ArtistFormModal from '@/features/artists/components/ArtistFormModal'
import ArtistCombobox from '@/features/artists/components/ArtistCombobox'
import { useCharacters } from '@/features/characters/api/useCharacters'
import { useFiles, FILES_QUERY_KEY } from '../api/useFiles'
import { useCreateArtpiece } from '@/features/artpieces/api/useCreateArtpiece'
import ArtpieceFilesInput from './ArtpieceFilesInput'

interface CharacterToken {
  id: string
  name: string
}

interface ArtpieceFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  initialFileIds?: string[]
  onSuccess?: () => void
}

export default function ArtpieceFormModal({ open, onOpenChange, initialFileIds, onSuccess }: ArtpieceFormModalProps) {
  const [title, setTitle] = useState('')
  const [notes, setNotes] = useState('')
  const [selectedArtist, setSelectedArtist] = useState<Artist | null>(null)
  const [selectedCharacters, setSelectedCharacters] = useState<CharacterToken[]>([])
  const [fileIds, setFileIds] = useState<string[]>([])
  const [hasUploading, setHasUploading] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [artistModalOpen, setArtistModalOpen] = useState(false)

  const charactersQuery = useCharacters()
  const { data: inboxFiles = [] } = useFiles()
  const createArtpiece = useCreateArtpiece()
  const queryClient = useQueryClient()

  const allCharacters = charactersQuery.data ?? []

  const availableCharacters = allCharacters
    .map((c) => ({ id: c.id, name: c.name }))
    .filter((c) => !selectedCharacters.some((sc) => sc.id === c.id))

  const initialFilesForInput = initialFileIds?.map((id) => {
    const match = inboxFiles.find((f) => f.id === id)
    return { id, name: match?.name ?? null, notes: match?.notes ?? null }
  })

  function resetForm() {
    setTitle('')
    setNotes('')
    setSelectedArtist(null)
    setSelectedCharacters([])
    setFileIds([])
    setHasUploading(false)
  }

  useEffect(() => {
    if (!open) resetForm()
  }, [open])

  function handleOpenChange(next: boolean) {
    if (!next && hasUploading) return
    if (!next) resetForm()
    onOpenChange(next)
  }

  const handleFileIdsChange = useCallback((ids: string[]) => setFileIds(ids), [])
  const handleHasUploading = useCallback((uploading: boolean) => setHasUploading(uploading), [])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!title.trim() || hasUploading) return
    setIsSubmitting(true)
    try {
      await createArtpiece.mutateAsync({
        title: title.trim(),
        notes: notes.trim() || null,
        artistId: selectedArtist?.id ?? null,
        characterIds: selectedCharacters.map((c) => c.id),
        fileIds,
      })
      await queryClient.invalidateQueries({ queryKey: FILES_QUERY_KEY })
      toast.success('Artpiece created')
      onSuccess?.()
      handleOpenChange(false)
    } catch {
      toast.error('Failed to create artpiece')
    } finally {
      setIsSubmitting(false)
    }
  }

  const canSave = title.trim().length > 0 && !hasUploading && !isSubmitting

  return (
    <>
      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent className="flex max-h-[calc(100dvh-2rem)] max-w-md flex-col overflow-hidden">
          <DialogHeader>
            <DialogTitle>New artpiece</DialogTitle>
          </DialogHeader>

          <form onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto pr-1">
            {/* Title */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="artpiece-title">
                Title <span className="text-destructive">*</span>
              </label>
              <Input
                id="artpiece-title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Artpiece title"
                disabled={isSubmitting}
              />
            </div>

            {/* Notes */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="artpiece-notes">
                Notes
              </label>
              <textarea
                id="artpiece-notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Notes about this artpiece..."
                rows={3}
                disabled={isSubmitting}
                className="flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>

            {/* Artist */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Artist</label>
              <div className="flex gap-2">
                <div className="flex-1">
                  <ArtistCombobox
                    value={selectedArtist}
                    onChange={setSelectedArtist}
                    disabled={isSubmitting}
                  />
                </div>
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  onClick={() => setArtistModalOpen(true)}
                  disabled={isSubmitting}
                >
                  <Plus className="size-4" />
                </Button>
              </div>
            </div>

            {/* Characters */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Characters</label>
              <TokenInput
                items={selectedCharacters}
                onChange={setSelectedCharacters}
                suggestions={availableCharacters}
                placeholder="Search characters…"
                disabled={isSubmitting}
              />
            </div>

            {/* Files */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Files</label>
              <ArtpieceFilesInput
                initialFiles={initialFilesForInput}
                onFileIdsChange={handleFileIdsChange}
                onHasUploading={handleHasUploading}
                disabled={isSubmitting}
              />
            </div>
          </form>

          <DialogFooter>
            <Button
              variant="outline"
              type="button"
              onClick={() => handleOpenChange(false)}
              disabled={isSubmitting || hasUploading}
            >
              Cancel
            </Button>
            <Button type="submit" onClick={handleSubmit} disabled={!canSave}>
              Save
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ArtistFormModal
        open={artistModalOpen}
        onOpenChange={setArtistModalOpen}
        onCreated={(artist) => setSelectedArtist(artist)}
      />
    </>
  )
}
