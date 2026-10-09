import { useState } from 'react'
import { Link } from 'react-router-dom'
import { RefreshCw, Upload, Plus, ExternalLink, Images, ReceiptText, type LucideIcon } from 'lucide-react'
import { useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip'
import { Loader2 } from 'lucide-react'
import AddFilesModal from '@/features/files/components/AddFilesModal'
import { useDashboard, DASHBOARD_QUERY_KEY, type DashboardCommission, type DashboardHousekeepingItem } from '@/features/dashboard/api/useDashboard'
import { usePatchCommission } from '@/features/commissions/api/usePatchCommission'
import { STATUS_LABELS } from '@/features/commissions/statusOptions'
import { formatRelativeTime } from '@/features/commissions/utils/relativeTime'

const STALE_DAYS = 14

function isStale(lastContactedAt: string | null): boolean {
  if (!lastContactedAt) return false
  const diff = Date.now() - new Date(lastContactedAt).getTime()
  return diff > STALE_DAYS * 24 * 60 * 60 * 1000
}

function daysElapsed(createdAt: string): number {
  return Math.floor((Date.now() - new Date(createdAt).getTime()) / (1000 * 60 * 60 * 24))
}

// --- Recent Artpieces ---

function RecentArtpiecesTile({ title, artistName, thumbnailUrl }: {
  title: string | null
  artistName: string | null
  thumbnailUrl: string | null
}) {
  return (
    <div className="flex flex-col gap-1">
      <div className="aspect-square w-full overflow-hidden rounded-lg bg-muted">
        {thumbnailUrl ? (
          <img
            src={thumbnailUrl}
            alt={title ?? ''}
            className="h-full w-full object-cover"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center text-muted-foreground">
            <Images className="size-6 opacity-40" />
          </div>
        )}
      </div>
      <p className="truncate text-sm font-medium leading-tight">{title ?? '—'}</p>
      {artistName && (
        <p className="truncate text-xs text-muted-foreground">{artistName}</p>
      )}
    </div>
  )
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
  const { data, isLoading, isError } = useDashboard()
  const queryClient = useQueryClient()
  const patchCommission = usePatchCommission()

  function handleStamp(id: string) {
    patchCommission.mutate(
      { id, lastContactedAt: new Date().toISOString() },
      { onSuccess: () => queryClient.invalidateQueries({ queryKey: DASHBOARD_QUERY_KEY }) },
    )
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
                    <RecentArtpiecesTile
                      key={a.id}
                      title={a.title}
                      artistName={a.artist_name}
                      thumbnailUrl={a.thumbnail_url}
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
    </div>
  )
}
