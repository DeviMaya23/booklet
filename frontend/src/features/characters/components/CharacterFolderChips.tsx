import { AlertTriangle, ExternalLink, Leaf, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import type { PublicFolder } from '../api/usePublicFolders'

const bookleafHost = () => import.meta.env.VITE_BOOKLEAF_HOST ?? ''

interface FolderChipsViewProps {
  folders: { id: string; name: string }[]
  /** The live Bookleaf list, when it has loaded; supplies share tokens. */
  liveFolders: PublicFolder[] | undefined
}

/** Read-only folder chips for the detail header; chips link out only when a token is known. */
export function FolderChipsView({ folders, liveFolders }: FolderChipsViewProps) {
  const tokens = new Map((liveFolders ?? []).map((f) => [f.id, f.token]))

  return (
    <div className="flex flex-wrap items-center gap-2 text-sm">
      <span className="font-medium text-muted-foreground">Folders</span>
      {folders.length === 0 ? (
        <span className="text-muted-foreground">No folders attached</span>
      ) : (
        folders.map((folder) => {
          const token = tokens.get(folder.id)
          const chip = 'inline-flex items-center gap-1.5 rounded-lg bg-muted px-2.5 py-1 font-medium'
          return token ? (
            <a
              key={folder.id}
              href={`${bookleafHost()}/share/${token}`}
              target="_blank"
              rel="noreferrer"
              className={`${chip} hover:bg-accent`}
            >
              <Leaf className="size-3.5 text-muted-foreground" />
              {folder.name}
            </a>
          ) : (
            <span key={folder.id} className={chip}>
              <Leaf className="size-3.5 text-muted-foreground" />
              {folder.name}
            </span>
          )
        })
      )}
    </div>
  )
}

interface FolderChipsEditProps {
  selectedFolders: PublicFolder[]
  invalidFolders: { id: string; name: string }[]
  onRemove: (id: string) => void
  onDismissInvalid: (id: string) => void
  onChoose: () => void
  disabled?: boolean
}

/** Editable folder chips: remove control, Bookleaf link, stale-folder warnings, Choose folders button. */
export function FolderChipsEdit({
  selectedFolders,
  invalidFolders,
  onRemove,
  onDismissInvalid,
  onChoose,
  disabled,
}: FolderChipsEditProps) {
  return (
    <section className="flex flex-col gap-3" aria-label="Folders">
      <div className="flex flex-wrap items-center gap-2">
        <h2 className="text-lg font-bold">Folders</h2>
        <span
          className="rounded-md bg-muted px-1.5 py-0.5 text-xs text-muted-foreground"
          aria-label={`${selectedFolders.length} folders`}
        >
          {selectedFolders.length}
        </span>
        <p className="text-sm text-muted-foreground">
          Bookleaf folders whose images belong to this character.
        </p>
      </div>

      {(selectedFolders.length > 0 || invalidFolders.length > 0) && (
        <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
          {selectedFolders.map((folder) => (
            <div
              key={folder.id}
              className="flex items-center gap-2 rounded-lg border bg-background px-3 py-2 text-sm"
            >
              <Leaf className="size-4 shrink-0 text-muted-foreground" />
              <span className="flex-1 truncate font-medium">{folder.name}</span>
              {folder.token && (
                <a
                  href={`${bookleafHost()}/share/${folder.token}`}
                  target="_blank"
                  rel="noreferrer"
                  className="text-muted-foreground hover:text-foreground"
                  aria-label={`Open ${folder.name} in Bookleaf`}
                >
                  <ExternalLink className="size-4" />
                </a>
              )}
              <button
                type="button"
                className="text-muted-foreground hover:text-foreground"
                onClick={() => onRemove(folder.id)}
                disabled={disabled}
                aria-label={`Remove ${folder.name}`}
              >
                <X className="size-4" />
              </button>
            </div>
          ))}

          {invalidFolders.map((folder) => (
            <div
              key={folder.id}
              className="flex items-center gap-2 rounded-lg border border-dashed border-orange-300 bg-orange-50 px-3 py-2 text-sm dark:border-orange-900 dark:bg-orange-950/30"
            >
              <AlertTriangle className="size-4 shrink-0 text-orange-500" />
              <div className="min-w-0 flex-1">
                <p className="truncate font-semibold">{folder.name}</p>
                <p className="text-xs text-orange-800 dark:text-orange-300">
                  No longer public in Bookleaf.
                </p>
              </div>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => onDismissInvalid(folder.id)}
                disabled={disabled}
              >
                Remove
              </Button>
            </div>
          ))}
        </div>
      )}

      <Button
        type="button"
        variant="secondary"
        className="self-start"
        onClick={onChoose}
        disabled={disabled}
      >
        Choose folders
      </Button>
    </section>
  )
}
