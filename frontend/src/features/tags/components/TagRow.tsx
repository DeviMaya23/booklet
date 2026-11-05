import { useState, type KeyboardEvent } from 'react'
import { Link } from 'react-router-dom'
import { Check, Pencil, Trash2, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import type { Tag } from '../api/useTags'

interface TagRowProps {
  tag: Tag
  editing: boolean
  error?: string | null
  isSaving?: boolean
  onEdit: () => void
  onCancelEdit: () => void
  onSubmitName: (name: string) => void
  onDelete: () => void
}

export default function TagRow({
  tag,
  editing,
  error,
  isSaving = false,
  onEdit,
  onCancelEdit,
  onSubmitName,
  onDelete,
}: TagRowProps) {
  const count = tag.artpiece_count ?? 0
  const usageLabel = `${count} ${count === 1 ? 'artpiece' : 'artpieces'}`

  if (editing) {
    return (
      <li className="flex flex-col gap-1 rounded-lg border px-3 py-2">
        <EditForm
          initialName={tag.name}
          isSaving={isSaving}
          onSubmit={onSubmitName}
          onCancel={onCancelEdit}
        />
        {error && <p role="alert" className="text-sm text-destructive">{error}</p>}
      </li>
    )
  }

  return (
    <li className="flex items-center gap-3 rounded-lg border px-3 py-2">
      <span className="font-semibold">{tag.name}</span>
      {count === 0 && (
        <span className="rounded-md bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground">Unused</span>
      )}
      <div className="flex-1" />
      {count > 0 ? (
        <Link
          to={`/app/artpieces?tags=${tag.id}`}
          className="text-sm text-muted-foreground hover:text-foreground"
        >
          {usageLabel} →
        </Link>
      ) : (
        <span className="text-sm text-muted-foreground">{usageLabel}</span>
      )}
      <Button variant="ghost" size="icon-sm" onClick={onEdit} aria-label={`Edit ${tag.name}`}>
        <Pencil className="size-4" />
      </Button>
      <Button variant="ghost" size="icon-sm" onClick={onDelete} aria-label={`Delete ${tag.name}`}>
        <Trash2 className="size-4" />
      </Button>
    </li>
  )
}

interface EditFormProps {
  initialName: string
  isSaving: boolean
  onSubmit: (name: string) => void
  onCancel: () => void
}

// Mounted only while a row is being edited, so the draft starts fresh from the
// tag's current name every time editing begins.
function EditForm({ initialName, isSaving, onSubmit, onCancel }: EditFormProps) {
  const [draft, setDraft] = useState(initialName)

  function submit() {
    const trimmed = draft.trim()
    if (trimmed === '') return
    if (trimmed === initialName) {
      onCancel()
      return
    }
    onSubmit(trimmed)
  }

  function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Enter') {
      e.preventDefault()
      submit()
    } else if (e.key === 'Escape') {
      e.preventDefault()
      onCancel()
    }
  }

  return (
    <div className="flex items-center gap-2">
      <Input
        autoFocus
        aria-label="Tag name"
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        onKeyDown={handleKeyDown}
        disabled={isSaving}
        className="max-w-xs"
      />
      <Button size="icon-sm" onClick={submit} disabled={isSaving} aria-label="Save">
        <Check className="size-4" />
      </Button>
      <Button variant="ghost" size="icon-sm" onClick={onCancel} disabled={isSaving} aria-label="Cancel">
        <X className="size-4" />
      </Button>
      <span className="text-xs text-muted-foreground">Enter to save · Esc to cancel</span>
    </div>
  )
}
