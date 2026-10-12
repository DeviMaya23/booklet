import { useState } from 'react'
import { ExternalLink, Leaf, RefreshCw, X } from 'lucide-react'
import { Button, buttonVariants } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import type { PublicFolder } from '../api/usePublicFolders'
import { FolderPreview } from './FolderPreview'

interface ChooseFoldersPanelProps {
  selectedIds: Set<string>
  onAdd: (folder: PublicFolder) => void
  onRemove: (id: string) => void
  onClose: () => void
  folders: PublicFolder[] | null
  isLoading: boolean
  isError: boolean
  onRefresh: () => void
}

export function ChooseFoldersPanel({
  selectedIds,
  onAdd,
  onRemove,
  onClose,
  folders,
  isLoading,
  isError,
  onRefresh,
}: ChooseFoldersPanelProps) {
  const [search, setSearch] = useState('')
  const [activePreviewId, setActivePreviewId] = useState<string | null>(null)

  const bookleafHost = import.meta.env.VITE_BOOKLEAF_HOST ?? ''

  const filtered = (folders ?? []).filter((f) =>
    f.name.toLowerCase().includes(search.toLowerCase()),
  )

  function handleToggle(folder: PublicFolder) {
    if (selectedIds.has(folder.id)) {
      onRemove(folder.id)
    } else {
      onAdd(folder)
    }
  }

  function handlePreviewToggle(id: string) {
    setActivePreviewId((prev) => (prev === id ? null : id))
  }

  return (
    <div className="flex h-full flex-col gap-3">
      <div className="flex items-center justify-between">
        <span className="text-sm font-semibold">Choose folders</span>
        <button
          type="button"
          className="text-muted-foreground hover:text-foreground"
          onClick={onClose}
          aria-label="Close folder panel"
        >
          <X className="size-4" />
        </button>
      </div>

      <Input
        placeholder="Search folders…"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="h-8"
      />

      {/* Error state — no data loaded */}
      {isError && (
        <div className="flex flex-col gap-3 rounded-md border border-destructive/40 bg-destructive/5 p-4">
          <div className="flex items-start gap-3">
            <Leaf className="mt-0.5 size-4 shrink-0 text-destructive" />
            <div className="flex-1 text-sm">
              <p className="font-semibold text-destructive">Couldn't reach Bookleaf</p>
              <p className="text-muted-foreground">
                Your attached folders are unchanged. Previews and adding folders need a connection.
              </p>
            </div>
            <Button variant="outline" size="sm" onClick={onRefresh} className="shrink-0 gap-1.5">
              <RefreshCw className="size-3.5" />
              Retry
            </Button>
          </div>
        </div>
      )}

      {/* Non-error content area */}
      {!isError && (
        <>
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <div className="flex items-center gap-1">
              <Leaf className="size-3" />
              <span>Public folders from Bookleaf</span>
            </div>
            <button
              type="button"
              className="hover:text-foreground"
              onClick={onRefresh}
              aria-label="Refresh folder list"
            >
              <RefreshCw className="size-3.5" />
            </button>
          </div>

          {/* Empty state */}
          {!isLoading && folders !== null && folders.length === 0 && (
            <div className="flex flex-1 flex-col items-center justify-center gap-4 rounded-md border border-dashed p-8 text-center">
              <div className="flex flex-col items-center gap-2">
                <Leaf className="size-8 text-muted-foreground" />
                <p className="font-semibold">No public folders yet</p>
                <p className="text-sm text-muted-foreground">
                  Make a folder public in Bookleaf, then come back and refresh.
                </p>
              </div>
              <div className="flex gap-2">
                <a
                  href={bookleafHost}
                  target="_blank"
                  rel="noreferrer"
                  className={cn(buttonVariants({ size: 'sm' }), 'gap-1.5')}
                >
                  <ExternalLink className="size-3.5" />
                  Open Bookleaf
                </a>
                <Button variant="outline" size="sm" onClick={onRefresh} className="gap-1.5">
                  <RefreshCw className="size-3.5" />
                  Refresh
                </Button>
              </div>
            </div>
          )}

          {/* Loading skeleton */}
          {isLoading && (
            <div className="flex flex-col gap-1">
              {Array.from({ length: 3 }).map((_, i) => (
                <div key={i} className="h-10 animate-pulse rounded-md border bg-muted" />
              ))}
            </div>
          )}

          {/* Folder list */}
          {!isLoading && folders !== null && folders.length > 0 && (
            <div className="flex flex-col gap-1 overflow-y-auto">
              {filtered.map((folder) => {
                const isSelected = selectedIds.has(folder.id)
                const isPreviewOpen = activePreviewId === folder.id
                return (
                  <div
                    key={folder.id}
                    className="rounded-md border"
                  >
                    <div className="flex items-center gap-3 px-3 py-2">
                      <Checkbox
                        id={`folder-check-${folder.id}`}
                        checked={isSelected}
                        onCheckedChange={() => handleToggle(folder)}
                        aria-label={folder.name}
                      />
                      <label
                        htmlFor={`folder-check-${folder.id}`}
                        className="flex-1 cursor-pointer truncate text-sm font-medium"
                      >
                        {folder.name}
                      </label>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="shrink-0 gap-1"
                        onClick={() => handlePreviewToggle(folder.id)}
                      >
                        {isPreviewOpen ? (
                          <>
                            <span className="text-xs">↑</span> Hide
                          </>
                        ) : (
                          'Preview'
                        )}
                      </Button>
                      <a
                        href={`${bookleafHost}/share/${folder.token}`}
                        target="_blank"
                        rel="noreferrer"
                        className="shrink-0 text-muted-foreground hover:text-foreground"
                        aria-label={`Open ${folder.name} in Bookleaf`}
                      >
                        <ExternalLink className="size-3.5" />
                      </a>
                    </div>

                    {isPreviewOpen && (
                      <div className="border-t px-3 pb-3 pt-2">
                        <FolderPreview folderID={folder.id} />
                      </div>
                    )}
                  </div>
                )
              })}
            </div>
          )}
        </>
      )}
    </div>
  )
}
