import { toast } from 'sonner'
import TokenInput from '@/components/TokenInput'
import { useTags, type Tag } from '../api/useTags'
import { useCreateTag, TagConflictError } from '../api/useCreateTag'

interface TagPickerProps {
  value: Tag[]
  onChange: (tags: Tag[]) => void
  disabled?: boolean
  placeholder?: string
}

function findByName(tags: Tag[], name: string): Tag | undefined {
  const lower = name.toLowerCase()
  return tags.find((t) => t.name.toLowerCase() === lower)
}

export default function TagPicker({ value, onChange, disabled, placeholder = 'Add tag…' }: TagPickerProps) {
  const tagsQuery = useTags({ sort: 'recent' })
  const createTag = useCreateTag()
  const allTags = tagsQuery.data ?? []

  async function handleCreate(name: string): Promise<Tag | null> {
    const existing = findByName(allTags, name)
    if (existing) return existing
    try {
      return await createTag.mutateAsync(name)
    } catch (err) {
      if (err instanceof TagConflictError) {
        const { data } = await tagsQuery.refetch()
        const found = findByName(data ?? [], name)
        if (found) return found
      }
      toast.error('Failed to create tag')
      return null
    }
  }

  // Until the list has loaded, "no existing tag with this name" cannot be told
  // apart from "not loaded yet", so adding and creating are held back.
  const unavailable = tagsQuery.isPending || tagsQuery.isError

  return (
    <div className="flex flex-col gap-1">
      <TokenInput
        items={value}
        onChange={onChange}
        suggestions={allTags}
        placeholder={tagsQuery.isPending ? 'Loading…' : placeholder}
        disabled={disabled || unavailable}
        openOnFocus
        onCreate={handleCreate}
        canCreate={(raw) => !findByName(allTags, raw)}
      />
      {tagsQuery.isError && <p className="text-sm text-destructive">Failed to load tags.</p>}
    </div>
  )
}
