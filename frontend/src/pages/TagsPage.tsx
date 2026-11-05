import { useState } from 'react'
import { ChevronDown } from 'lucide-react'
import { Input } from '@/components/ui/input'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useTags } from '@/features/tags/api/useTags'
import { useTagManagement } from '@/features/tags/hooks/useTagManagement'
import TagsList from '@/features/tags/components/TagsList'
import DeleteTagDialog from '@/features/tags/components/DeleteTagDialog'
import MergeTagDialog from '@/features/tags/components/MergeTagDialog'

type TagSort = 'recent' | 'name'

const SORT_LABELS: Record<TagSort, string> = {
  recent: 'Recently used',
  name: 'Name (A–Z)',
}

export default function TagsPage() {
  const [search, setSearch] = useState('')
  const [sort, setSort] = useState<TagSort>('recent')

  const { data, isPending, isError } = useTags({
    sort: sort === 'recent' ? 'recent' : undefined,
    withUsage: true,
  })
  const tags = data ?? []
  const management = useTagManagement(tags)

  const term = search.trim().toLowerCase()
  const filtered = term ? tags.filter((t) => t.name.toLowerCase().includes(term)) : tags

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Input
          placeholder="Search tags..."
          aria-label="Search tags"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
        <div className="flex-1" />
        <DropdownMenu>
          <DropdownMenuTrigger
            aria-label="Sort tags"
            className="inline-flex h-9 items-center gap-1.5 rounded-md border border-input bg-background px-3 text-sm shadow-xs hover:bg-accent hover:text-accent-foreground focus-visible:outline-none"
          >
            {SORT_LABELS[sort]}
            <ChevronDown className="size-4" />
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onClick={() => setSort('recent')}>{SORT_LABELS.recent}</DropdownMenuItem>
            <DropdownMenuItem onClick={() => setSort('name')}>{SORT_LABELS.name}</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {isError && <p className="text-sm text-destructive">Failed to load tags.</p>}
      {!isPending && !isError && (
        <TagsList
          tags={filtered}
          isEmpty={tags.length === 0}
          editingId={management.editingId}
          renameError={management.renameError}
          isSaving={management.isRenaming}
          onEdit={management.startEdit}
          onCancelEdit={management.cancelEdit}
          onSubmitName={management.submitRename}
          onDelete={management.requestDelete}
        />
      )}

      <DeleteTagDialog
        item={management.deleteTarget}
        onCancel={management.cancelDelete}
        onConfirm={management.confirmDelete}
        isPending={management.isDeleting}
      />
      <MergeTagDialog
        item={management.mergeTarget}
        onCancel={management.cancelMerge}
        onConfirm={management.confirmMerge}
        isPending={management.isMerging}
      />
    </div>
  )
}
