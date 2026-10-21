import { useState } from 'react'
import { Link } from 'react-router-dom'
import { RefreshCw, Upload, Plus, ExternalLink, Images, ReceiptText, Trash2, ImageIcon, type LucideIcon } from 'lucide-react'
import { useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip'
import { Loader2 } from 'lucide-react'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Checkbox } from '@/components/ui/checkbox'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import AddFilesModal from '@/features/files/components/AddFilesModal'
import { useDashboard, DASHBOARD_QUERY_KEY, type DashboardCommission, type DashboardHousekeepingItem } from '@/features/dashboard/api/useDashboard'
import { usePatchCommission } from '@/features/commissions/api/usePatchCommission'
import { STATUS_LABELS } from '@/features/commissions/statusOptions'
import { formatRelativeTime } from '@/features/commissions/utils/relativeTime'
import ThumbnailTile, { type ThumbnailTileMenuItem } from '@/components/ThumbnailTile'
import { useSavedFilters, SAVED_FILTERS_QUERY_KEY } from '@/features/artpieces/api/useSavedFilters'
import { useDeleteSavedFilter } from '@/features/artpieces/api/useDeleteSavedFilter'
import { parseFilterPayload } from '@/features/artpieces/lib/parseFilterPayload'
import ChangeThumbnailModal from '@/features/artpieces/components/ChangeThumbnailModal'
import { type SavedFilter } from '@/features/artpieces/api/useCreateSavedFilter'

const STALE_DAYS = 14

function isStale(lastContactedAt: string | null): boolean {
  if (!lastContactedAt) return false
  const diff = Date.now() - new Date(lastContactedAt).getTime()
  return diff > STALE_DAYS * 24 * 60 * 60 * 1000
}

function daysElapsed(createdAt: string): number {
  return Math.floor((Date.now() - new Date(createdAt).getTime()) / (1000 * 60 * 60 * 24))
}

// --- Helpers ---

function filterToHref(filter: SavedFilter): string {
  const payload = parseFilterPayload(filter.filter_payload)
  if (!payload) return '/app/artpieces'
  const params = new URLSearchParams()
  if (payload.artist_id) params.set('artist', payload.artist_id)
  if (payload.character_ids.length > 0) params.set('characters', payload.character_ids.join(','))
  if (payload.match_mode === 'any') params.set('match', 'any')
  const qs = params.toString()
  return qs ? `/app/artpieces?${qs}` : '/app/artpieces'
}

function getHiddenFilterIds(): string[] {
  try {
    const raw = localStorage.getItem('dashboard.hiddenFilterIds')
    if (!raw) return []
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed.filter((id): id is string => typeof id === 'string') : []
  } catch {
    return []
  }
}

function setHiddenFilterIds(ids: string[]): void {
  try {
    localStorage.setItem('dashboard.hiddenFilterIds', JSON.stringify(ids))
  } catch {
    // ignore
  }
}

// --- In Progress Row ---

function InProgressRow({ commission, onStamp }: {
  commission: DashboardCommission
  onStamp: (id: string) => void
}) {
  const stale = isStale(commission.last_contacted_at)
  const days = daysElapsed(commission.created_at)

  return (
    <div className="flex items-center gap-3 py-2 text-sm">
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-1.5">
          <Link
            to="/app/commissions"
            className="truncate font-medium hover:underline"
          >
            {commission.title ?? '—'}
          </Link>
          {commission.artist_link && (
            <a
              href={commission.artist_link}
              target="_blank"
              rel="noopener noreferrer"
              className="shrink-0 text-muted-foreground hover:text-foreground"
            >
              <ExternalLink className="size-3" />
            </a>
          )}
        </div>
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          {commission.artist_name && <span>{commission.artist_name}</span>}
          <span>·</span>
          <span>{STATUS_LABELS[commission.status]}</span>
          <span>·</span>
          <span>{commission.paid ? 'Paid' : 'Unpaid'}</span>
          <span>·</span>
          <span>waiting {days}d</span>
          <span>·</span>
          <span className={stale ? 'text-amber-600 dark:text-amber-400' : ''}>
            {commission.last_contacted_at
              ? formatRelativeTime(new Date(commission.last_contacted_at))
              : 'not contacted yet'}
          </span>
        </div>
      </div>
      <Tooltip>
        <TooltipTrigger
          onClick={() => onStamp(commission.id)}
          aria-label="Mark contacted now"
          className="shrink-0 inline-flex size-7 items-center justify-center rounded text-muted-foreground hover:bg-accent hover:text-foreground"
        >
          <RefreshCw className="size-3.5" />
        </TooltipTrigger>
        <TooltipContent side="top">Mark as contacted today</TooltipContent>
      </Tooltip>
    </div>
  )
}

// --- Housekeeping ---

function HousekeepingSection({ label, items, icon: Icon, linkTo }: {
  label: string
  items: DashboardHousekeepingItem[]
  icon: LucideIcon
  linkTo: (item: DashboardHousekeepingItem) => string
}) {
  if (items.length === 0) return null
  return (
    <div>
      <p className="mb-1 text-xs font-medium text-muted-foreground">{label}</p>
      <ul className="space-y-0.5">
        {items.map((item) => (
          <li key={item.id}>
            <Link
              to={linkTo(item)}
              className="flex items-center justify-between rounded px-2 py-1.5 text-sm hover:bg-accent"
            >
              <span className="flex items-center gap-2">
                <Icon className="size-3.5 shrink-0 text-muted-foreground" />
                <span className="truncate">{item.title ?? '—'}</span>
              </span>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  )
}

// --- Page ---

export default function DashboardPage() {
  const [addFilesOpen, setAddFilesOpen] = useState(false)
  const [deleteTargetId, setDeleteTargetId] = useState<string | null>(null)
  const [changeThumbnailFilter, setChangeThumbnailFilter] = useState<SavedFilter | null>(null)
  const [hiddenIds, setHiddenIdsState] = useState<string[]>(() => getHiddenFilterIds())
  const [chooseShownOpen, setChooseShownOpen] = useState(false)

  const { data, isLoading, isError } = useDashboard()
  const savedFiltersQuery = useSavedFilters()
  const queryClient = useQueryClient()
  const patchCommission = usePatchCommission()
  const deleteFilter = useDeleteSavedFilter()

  function handleStamp(id: string) {
    patchCommission.mutate(
      { id, lastContactedAt: new Date().toISOString() },
      { onSuccess: () => queryClient.invalidateQueries({ queryKey: DASHBOARD_QUERY_KEY }) },
    )
  }

  function handleDeleteConfirm() {
    if (!deleteTargetId) return
    deleteFilter.mutate(deleteTargetId, {
      onSuccess: () => queryClient.invalidateQueries({ queryKey: SAVED_FILTERS_QUERY_KEY }),
      onSettled: () => setDeleteTargetId(null),
    })
  }

  function toggleHidden(id: string) {
    setHiddenIdsState((prev) => {
      const next = prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
      setHiddenFilterIds(next)
      return next
    })
  }

  const housekeepingItems = data ? [
    data.housekeeping.commissions_no_artist,
    data.housekeeping.artpieces_no_artist,
    data.housekeeping.done_no_artpieces,
    data.housekeeping.artpieces_no_files,
  ] : []
  const totalHousekeeping = housekeepingItems.reduce((sum, arr) => sum + arr.length, 0)

  return (
    <div className="flex h-full flex-col gap-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={() => setAddFilesOpen(true)}>
            <Upload className="mr-1.5 size-4" />
            Upload
          </Button>
          <Button size="sm" disabled>
            <Plus className="mr-1.5 size-4" />
            New commission
          </Button>
        </div>
      </div>

      {isLoading && (
        <div className="flex flex-1 items-center justify-center">
          <Loader2 className="size-6 animate-spin text-muted-foreground" />
        </div>
      )}

      {isError && (
        <div className="flex flex-1 items-center justify-center">
          <p className="text-sm text-muted-foreground">Failed to load dashboard. Please try again.</p>
        </div>
      )}

      {data && (
        <div className="flex min-h-0 flex-1 gap-6">
          {/* Left column */}
          <div className="flex min-w-0 flex-1 flex-col gap-6">
            {/* Saved filters */}
            {(savedFiltersQuery.data !== undefined || savedFiltersQuery.isError) && (
              <section>
                <div className="mb-3 flex items-center justify-between">
                  <h2 className="text-base font-semibold">Saved filters</h2>
                  {savedFiltersQuery.data != null && savedFiltersQuery.data.length > 0 && (
                    <Popover open={chooseShownOpen} onOpenChange={setChooseShownOpen}>
                      <PopoverTrigger className="text-sm text-muted-foreground hover:text-foreground">
                        Choose shown
                      </PopoverTrigger>
                      <PopoverContent align="end" className="w-56 p-2">
                        <p className="mb-1.5 px-1 text-xs font-medium text-muted-foreground">Show on dashboard</p>
                        {savedFiltersQuery.data.map((sf) => (
                          <label key={sf.id} className="flex cursor-pointer items-center gap-2 rounded px-1 py-1.5 hover:bg-accent">
                            <Checkbox
                              checked={!hiddenIds.includes(sf.id)}
                              onCheckedChange={() => toggleHidden(sf.id)}
                            />
                            <span className="truncate text-sm">{sf.name}</span>
                          </label>
                        ))}
                      </PopoverContent>
                    </Popover>
                  )}
                </div>
                {savedFiltersQuery.isError ? (
                  <p className="text-sm text-muted-foreground">Failed to load saved filters.</p>
                ) : savedFiltersQuery.data!.length === 0 ? (
                  <p className="text-sm text-muted-foreground">No saved filters yet. Save a filter from the artpieces page to show it here.</p>
                ) : savedFiltersQuery.data.every((sf) => hiddenIds.includes(sf.id)) ? (
                  <p className="text-sm text-muted-foreground">All filters are hidden. Use "Choose shown" to show some.</p>
                ) : (
                  <div className="grid grid-cols-5 gap-3">
                    {savedFiltersQuery.data
                      .filter((sf) => !hiddenIds.includes(sf.id))
                      .map((sf) => {
                        const menuItems: ThumbnailTileMenuItem[] = [
                          {
                            label: 'Change thumbnail',
                            icon: ImageIcon,
                            onClick: () => setChangeThumbnailFilter(sf),
                          },
                          {
                            label: 'Delete',
                            icon: Trash2,
                            variant: 'destructive',
                            onClick: () => setDeleteTargetId(sf.id),
                          },
                        ]
                        return (
                          <ThumbnailTile
                            key={sf.id}
                            imageUrl={sf.thumbnail_url}
                            label={sf.name}
                            href={filterToHref(sf)}
                            menuItems={menuItems}
                          />
                        )
                      })}
                  </div>
                )}
              </section>
            )}

            {/* Recent artpieces */}
            <section>
              <div className="mb-3 flex items-center justify-between">
                <h2 className="text-base font-semibold">Recent artpieces</h2>
                <Link to="/app/artpieces" className="text-sm text-muted-foreground hover:text-foreground">
                  View all →
                </Link>
              </div>
              {data.recent_artpieces.length === 0 ? (
                <p className="text-sm text-muted-foreground">No artpieces yet.</p>
              ) : (
                <div className="grid grid-cols-5 gap-3">
                  {data.recent_artpieces.map((a) => (
                    <ThumbnailTile
                      key={a.id}
                      imageUrl={a.thumbnail_url}
                      label={a.title ?? '—'}
                      sublabel={a.artist_name ?? undefined}
                      href={`/app/artpieces/${a.id}`}
                    />
                  ))}
                </div>
              )}
            </section>
          </div>

          {/* Right column */}
          <div className="flex w-80 shrink-0 flex-col gap-4">
            {/* In progress */}
            <section className="rounded-xl border bg-card p-4">
              <div className="mb-3 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <h2 className="text-sm font-semibold">In progress</h2>
                  <span className="rounded-full bg-muted px-1.5 py-0.5 text-xs font-medium text-muted-foreground">
                    {data.in_progress.length}
                  </span>
                </div>
                <Link to="/app/commissions" className="text-xs text-muted-foreground hover:text-foreground">
                  View all →
                </Link>
              </div>
              {data.in_progress.length === 0 ? (
                <p className="text-sm text-muted-foreground">Nothing in progress.</p>
              ) : (
                <div className="divide-y">
                  {data.in_progress.map((c) => (
                    <InProgressRow key={c.id} commission={c} onStamp={handleStamp} />
                  ))}
                </div>
              )}
            </section>

            {/* Housekeeping */}
            <section className="rounded-xl border bg-card p-4">
              <div className="mb-3 flex items-center gap-2">
                <h2 className="text-sm font-semibold">Housekeeping</h2>
                {totalHousekeeping > 0 && (
                  <span className="rounded-full bg-muted px-1.5 py-0.5 text-xs font-medium text-muted-foreground">
                    {totalHousekeeping}
                  </span>
                )}
              </div>
              {totalHousekeeping === 0 ? (
                <p className="text-sm text-muted-foreground">All tidy. Nothing needs fixing.</p>
              ) : (
                <div className="flex flex-col gap-4">
                  <HousekeepingSection
                    label="No artist"
                    items={data.housekeeping.commissions_no_artist}
                    icon={ReceiptText}
                    linkTo={() => '/app/commissions'}
                  />
                  <HousekeepingSection
                    label="No artist"
                    items={data.housekeeping.artpieces_no_artist}
                    icon={Images}
                    linkTo={(item) => `/app/artpieces/${item.id}`}
                  />
                  <HousekeepingSection
                    label="Done, no artpieces"
                    items={data.housekeeping.done_no_artpieces}
                    icon={ReceiptText}
                    linkTo={() => '/app/commissions'}
                  />
                  <HousekeepingSection
                    label="No files"
                    items={data.housekeeping.artpieces_no_files}
                    icon={Images}
                    linkTo={(item) => `/app/artpieces/${item.id}`}
                  />
                </div>
              )}
            </section>
          </div>
        </div>
      )}

      <AddFilesModal open={addFilesOpen} onOpenChange={setAddFilesOpen} />

      <AlertDialog open={deleteTargetId != null} onOpenChange={(open) => { if (!open) setDeleteTargetId(null) }}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete saved filter?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently remove the saved filter. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteConfirm}>Delete</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {changeThumbnailFilter && (
        <ChangeThumbnailModal
          filter={changeThumbnailFilter}
          open={true}
          onOpenChange={(open) => { if (!open) setChangeThumbnailFilter(null) }}
          onSaved={() => {
            queryClient.invalidateQueries({ queryKey: SAVED_FILTERS_QUERY_KEY })
            setChangeThumbnailFilter(null)
          }}
        />
      )}
    </div>
  )
}
