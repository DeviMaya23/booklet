import { useState, useRef, useEffect } from 'react'
import { toast } from 'sonner'
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from '@dnd-kit/core'
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { GripVertical, Pencil, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { type Verse } from '../api/useVerses'
import { useCreateVerse } from '../api/useCreateVerse'
import { useUpdateVerse } from '../api/useUpdateVerse'
import { useDeleteVerse } from '../api/useDeleteVerse'
import { useReorderVerses } from '../api/useReorderVerses'
import { type Character } from '@/features/characters/api/useCharacters'

interface VersesModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  verses: Verse[]
  characters: Character[]
}

interface SortableVerseRowProps {
  verse: Verse
  characterCount: number
  onRenameStart: (id: string) => void
  onDeleteStart: (id: string) => void
}

function SortableVerseRow({
  verse,
  characterCount,
  onRenameStart,
  onDeleteStart,
}: SortableVerseRowProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: verse.id,
  })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="flex items-center gap-2 rounded-md px-2 py-1.5 hover:bg-muted/50 group"
    >
      <button
        type="button"
        className="cursor-grab text-muted-foreground hover:text-foreground touch-none"
        aria-label="Drag to reorder"
        {...attributes}
        {...listeners}
      >
        <GripVertical className="size-4" />
      </button>
      <span className="flex-1 text-sm">{verse.name}</span>
      <span className="text-xs text-muted-foreground mr-1">
        {characterCount} {characterCount === 1 ? 'character' : 'characters'}
      </span>
      <button
        type="button"
        className="text-muted-foreground hover:text-foreground opacity-0 group-hover:opacity-100 transition-opacity"
        onClick={() => onRenameStart(verse.id)}
        aria-label={`Rename ${verse.name}`}
      >
        <Pencil className="size-3.5" />
      </button>
      <button
        type="button"
        className="text-muted-foreground hover:text-destructive opacity-0 group-hover:opacity-100 transition-opacity"
        onClick={() => onDeleteStart(verse.id)}
        aria-label={`Delete ${verse.name}`}
      >
        <Trash2 className="size-3.5" />
      </button>
    </div>
  )
}

interface RenameRowProps {
  verse: Verse
  onSave: (id: string, name: string) => void
  onCancel: () => void
  error: string | null
}

function RenameRow({ verse, onSave, onCancel, error }: RenameRowProps) {
  const [value, setValue] = useState(verse.name)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    inputRef.current?.focus()
    inputRef.current?.select()
  }, [])

  return (
    <div className="flex flex-col gap-1 px-2 py-1.5">
      <div className="flex items-center gap-2">
        <div className="size-4 shrink-0" />
        <Input
          ref={inputRef}
          value={value}
          onChange={(e) => setValue(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              e.preventDefault()
              if (value.trim()) onSave(verse.id, value.trim())
            } else if (e.key === 'Escape') {
              onCancel()
            }
          }}
          className="h-7 text-sm flex-1"
        />
        <Button size="sm" variant="ghost" className="h-7 px-2 text-xs" onClick={() => { if (value.trim()) onSave(verse.id, value.trim()) }}>Save</Button>
        <Button size="sm" variant="ghost" className="h-7 px-2 text-xs text-muted-foreground" onClick={onCancel}>Cancel</Button>
      </div>
      {error && <p className="text-xs text-destructive pl-6">{error}</p>}
    </div>
  )
}

interface DeleteConfirmRowProps {
  verse: Verse
  characterCount: number
  onConfirm: (id: string) => void
  onCancel: () => void
  isPending: boolean
}

function DeleteConfirmRow({ verse, characterCount, onConfirm, onCancel, isPending }: DeleteConfirmRowProps) {
  return (
    <div className="flex flex-col gap-2 rounded-md bg-destructive/5 border border-destructive/20 px-3 py-2.5">
      <p className="text-sm">
        <span className="font-medium">"{verse.name}"</span>
        {' '}— {characterCount} {characterCount === 1 ? 'character' : 'characters'} will move to Unsorted.
      </p>
      <div className="flex gap-2">
        <Button
          size="sm"
          variant="destructive"
          className="h-7 text-xs"
          onClick={() => onConfirm(verse.id)}
          disabled={isPending}
        >
          Delete
        </Button>
        <Button
          size="sm"
          variant="outline"
          className="h-7 text-xs"
          onClick={onCancel}
          disabled={isPending}
        >
          Cancel
        </Button>
      </div>
    </div>
  )
}

interface NewVerseRowProps {
  onSave: (name: string) => void
  onCancel: () => void
  error: string | null
  isPending: boolean
}

function NewVerseRow({ onSave, onCancel, error, isPending }: NewVerseRowProps) {
  const [value, setValue] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  return (
    <div className="flex flex-col gap-1 px-2 py-1.5">
      <div className="flex items-center gap-2">
        <div className="size-4 shrink-0" />
        <Input
          ref={inputRef}
          value={value}
          onChange={(e) => setValue(e.target.value)}
          placeholder="Verse name"
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              e.preventDefault()
              if (value.trim()) onSave(value.trim())
            } else if (e.key === 'Escape') {
              onCancel()
            }
          }}
          className="h-7 text-sm flex-1"
          disabled={isPending}
        />
        <Button size="sm" variant="ghost" className="h-7 px-2 text-xs" onClick={() => { if (value.trim()) onSave(value.trim()) }} disabled={isPending || !value.trim()}>Save</Button>
        <Button size="sm" variant="ghost" className="h-7 px-2 text-xs text-muted-foreground" onClick={onCancel} disabled={isPending}>Cancel</Button>
      </div>
      {error && <p className="text-xs text-destructive pl-6">{error}</p>}
    </div>
  )
}

export default function VersesModal({ open, onOpenChange, verses, characters }: VersesModalProps) {
  const [order, setOrder] = useState<string[]>(() => verses.map((v) => v.id))
  const [renamingId, setRenamingId] = useState<string | null>(null)
  const [renameError, setRenameError] = useState<string | null>(null)
  const [deletingId, setDeletingId] = useState<string | null>(null)
  const [addingNew, setAddingNew] = useState(false)
  const [newVerseError, setNewVerseError] = useState<string | null>(null)

  const createMutation = useCreateVerse()
  const updateMutation = useUpdateVerse()
  const deleteMutation = useDeleteVerse()
  const reorderMutation = useReorderVerses()

  useEffect(() => {
    setOrder(verses.map((v) => v.id))
  }, [verses])

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  )

  const verseMap = new Map(verses.map((v) => [v.id, v]))
  const orderedVerses = order.map((id) => verseMap.get(id)).filter((v): v is Verse => v !== undefined)

  function charCountFor(verseId: string): number {
    return characters.filter((c) => c.verse_id === verseId).length
  }

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event
    if (!over || active.id === over.id) return
    const oldIndex = order.indexOf(active.id as string)
    const newIndex = order.indexOf(over.id as string)
    const previousOrder = order
    const newOrder = arrayMove(order, oldIndex, newIndex)
    setOrder(newOrder)
    reorderMutation.mutate(newOrder, {
      onError: () => {
        setOrder(previousOrder)
        toast.error('Failed to reorder verses')
      },
    })
  }

  function handleRenameStart(id: string) {
    setRenamingId(id)
    setRenameError(null)
    setDeletingId(null)
    setAddingNew(false)
  }

  function handleRenameSave(id: string, name: string) {
    updateMutation.mutate(
      { id, name },
      {
        onSuccess: () => {
          setRenamingId(null)
          setRenameError(null)
        },
        onError: (err) => {
          const status = (err as Error & { status?: number }).status
          setRenameError(status === 409 ? 'A verse with that name already exists.' : 'Failed to rename verse.')
        },
      },
    )
  }

  function handleDeleteStart(id: string) {
    setDeletingId(id)
    setRenamingId(null)
    setRenameError(null)
    setAddingNew(false)
  }

  function handleDeleteConfirm(id: string) {
    deleteMutation.mutate(id, {
      onSuccess: () => {
        setDeletingId(null)
      },
      onError: () => {
        setDeletingId(null)
      },
    })
  }

  function handleNewVerseSave(name: string) {
    createMutation.mutate(
      { name },
      {
        onSuccess: () => {
          setAddingNew(false)
          setNewVerseError(null)
        },
        onError: (err) => {
          const status = (err as Error & { status?: number }).status
          setNewVerseError(status === 409 ? 'A verse with that name already exists.' : 'Failed to create verse.')
        },
      },
    )
  }

  const showEmpty = orderedVerses.length === 0 && !addingNew

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Verses</DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-1 min-h-[120px]">
          {showEmpty ? (
            <div className="flex flex-col items-center justify-center gap-1 py-8 text-center">
              <p className="text-sm font-medium">No verses yet</p>
              <p className="text-xs text-muted-foreground">Create one to group characters on the list.</p>
            </div>
          ) : (
            <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
              <SortableContext items={order} strategy={verticalListSortingStrategy}>
                {orderedVerses.map((verse) => {
                  if (deletingId === verse.id) {
                    return (
                      <DeleteConfirmRow
                        key={verse.id}
                        verse={verse}
                        characterCount={charCountFor(verse.id)}
                        onConfirm={handleDeleteConfirm}
                        onCancel={() => setDeletingId(null)}
                        isPending={deleteMutation.isPending}
                      />
                    )
                  }
                  if (renamingId === verse.id) {
                    return (
                      <RenameRow
                        key={verse.id}
                        verse={verse}
                        onSave={handleRenameSave}
                        onCancel={() => { setRenamingId(null); setRenameError(null) }}
                        error={renameError}
                      />
                    )
                  }
                  return (
                    <SortableVerseRow
                      key={verse.id}
                      verse={verse}
                      characterCount={charCountFor(verse.id)}
                      onRenameStart={handleRenameStart}
                      onDeleteStart={handleDeleteStart}
                    />
                  )
                })}
              </SortableContext>
            </DndContext>
          )}

          {addingNew && (
            <NewVerseRow
              onSave={handleNewVerseSave}
              onCancel={() => { setAddingNew(false); setNewVerseError(null) }}
              error={newVerseError}
              isPending={createMutation.isPending}
            />
          )}
        </div>

        <div className="border-t pt-3">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => {
              setAddingNew(true)
              setRenamingId(null)
              setDeletingId(null)
            }}
            disabled={addingNew}
          >
            + New verse
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
