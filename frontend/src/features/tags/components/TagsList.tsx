import TagRow from './TagRow'
import type { Tag } from '../api/useTags'

interface TagsListProps {
  tags: Tag[]
  // True when the user owns no tags at all (as opposed to search matching none).
  isEmpty: boolean
  editingId: string | null
  renameError: string | null
  isSaving: boolean
  onEdit: (tag: Tag) => void
  onCancelEdit: () => void
  onSubmitName: (tag: Tag, name: string) => void
  onDelete: (tag: Tag) => void
}

export default function TagsList({
  tags,
  isEmpty,
  editingId,
  renameError,
  isSaving,
  onEdit,
  onCancelEdit,
  onSubmitName,
  onDelete,
}: TagsListProps) {
  if (isEmpty) {
    return (
      <div className="flex flex-col items-center gap-1 py-16 text-center">
        <p className="font-semibold">No tags yet</p>
        <p className="text-sm text-muted-foreground">Tags you add from an artpiece’s edit page will show up here.</p>
      </div>
    )
  }

  return (
    <ul className="flex flex-col gap-2">
      {tags.map((tag) => (
        <TagRow
          key={tag.id}
          tag={tag}
          editing={editingId === tag.id}
          error={editingId === tag.id ? renameError : null}
          isSaving={editingId === tag.id && isSaving}
          onEdit={() => onEdit(tag)}
          onCancelEdit={onCancelEdit}
          onSubmitName={(name) => onSubmitName(tag, name)}
          onDelete={() => onDelete(tag)}
        />
      ))}
    </ul>
  )
}
