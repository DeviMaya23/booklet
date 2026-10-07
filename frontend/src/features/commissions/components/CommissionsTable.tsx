import { ExternalLink, RefreshCw, ChevronUp, ChevronDown, AlertCircle, CalendarIcon, Pencil, Trash2, Info } from 'lucide-react'
import StatusChip from './StatusChip'
import ArtpieceThumbnailStack from './ArtpieceThumbnailStack'
import {
  Table, TableHeader, TableBody, TableHead, TableRow, TableCell,
} from '@/components/ui/table'
import { Checkbox } from '@/components/ui/checkbox'
import { Popover, PopoverTrigger, PopoverContent } from '@/components/ui/popover'
import { Calendar } from '@/components/ui/calendar'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider } from '@/components/ui/tooltip'
import { type Commission } from '../api/useCommissions'
import { usePatchCommission } from '../api/usePatchCommission'
import { formatRelativeTime } from '../utils/relativeTime'
import { useCommissionsSort, type SortKey, type SortDir } from '../hooks/useCommissionsSort'
import { useCellErrors } from '../hooks/useCellErrors'
import { useCommissionActions } from '../hooks/useCommissionActions'

const STALE_DAYS = 14

function isStale(lastContactedAt: string | null): boolean {
  if (!lastContactedAt) return false
  const diff = Date.now() - new Date(lastContactedAt).getTime()
  return diff > STALE_DAYS * 24 * 60 * 60 * 1000
}

function LastContactCell({ commission, onStamp, hasError }: {
  commission: Commission
  onStamp: () => void
  hasError: boolean
}) {
  const stale = isStale(commission.last_contacted_at)
  return (
    <span className={[
      'inline-flex items-center gap-1.5 rounded px-1.5 py-1 transition-colors',
      hasError ? 'ring-1 ring-destructive' : '',
    ].join(' ')}>
      <span className={[
        stale ? 'text-amber-600 dark:text-amber-400' : '',
        !commission.last_contacted_at ? 'text-muted-foreground' : '',
      ].join(' ')}>
        {commission.last_contacted_at
          ? formatRelativeTime(new Date(commission.last_contacted_at))
          : '—'
        }
      </span>
      <Tooltip>
        <TooltipTrigger
          onClick={onStamp}
          aria-label="Mark contacted now"
          className={[
            'inline-flex size-[1.375rem] items-center justify-center rounded hover:bg-accent',
            stale ? 'text-amber-600 dark:text-amber-400 hover:text-amber-600' : 'text-muted-foreground hover:text-foreground',
          ].join(' ')}
        >
          <RefreshCw className="size-3.5" />
        </TooltipTrigger>
        <TooltipContent side="top">Mark as contacted today</TooltipContent>
      </Tooltip>
      <CellError show={hasError} />
    </span>
  )
}

interface CommissionsTableProps {
  commissions: Commission[]
  view: 'active' | 'done'
  onEdit: (commission: Commission) => void
  onDelete: (commission: Commission) => void
  onCommissionDone?: (id: string, prevStatus: string, title: string | null) => void
}

function SortIcon({ column, sortKey, sortDir }: { column: SortKey; sortKey: SortKey; sortDir: SortDir }) {
  if (column !== sortKey) return null
  return sortDir === 'asc' ? <ChevronUp className="size-3" /> : <ChevronDown className="size-3" />
}

function SortableHead({ label, column, sortKey, sortDir, onToggle }: {
  label: string
  column: SortKey
  sortKey: SortKey
  sortDir: SortDir
  onToggle: (col: SortKey) => void
}) {
  return (
    <TableHead>
      <button
        type="button"
        onClick={() => onToggle(column)}
        className="inline-flex items-center gap-1 text-xs font-medium text-muted-foreground hover:text-foreground"
      >
        {label}
        <SortIcon column={column} sortKey={sortKey} sortDir={sortDir} />
      </button>
    </TableHead>
  )
}

interface CellErrorProps { show: boolean }
function CellError({ show }: CellErrorProps) {
  if (!show) return null
  return <AlertCircle className="size-3.5 text-destructive shrink-0" aria-label="Save failed" />
}

export default function CommissionsTable({ commissions, view, onEdit, onDelete, onCommissionDone }: CommissionsTableProps) {
  const { sorted, sortKey, sortDir, toggleSort } = useCommissionsSort(commissions)
  const { setCellError, hasCellError } = useCellErrors()
  const patchCommission = usePatchCommission()
  const { patchStatus, patchPaid, patchPaidDate, stampLastContacted } = useCommissionActions({
    mutateAsync: patchCommission.mutateAsync,
    setCellError,
    onStatusDone: view === 'active' ? onCommissionDone : undefined,
  })

  return (
    <TooltipProvider>
      <Table>
        <TableHeader>
          <TableRow>
            <SortableHead label="Title" column="title" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
            <SortableHead label="Artist" column="artist_name" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
            <SortableHead label="Status" column="status" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
            <SortableHead label="Paid" column="paid" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
            <SortableHead label="Paid date" column="paid_date" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
            {view === 'active' ? (
              <TableHead>
                <span className="inline-flex items-center gap-1 text-xs font-medium text-muted-foreground">
                  Last contact
                  <Tooltip>
                    <TooltipTrigger className="cursor-default inline-flex items-center">
                      <Info className="size-3 text-muted-foreground" />
                    </TooltipTrigger>
                    <TooltipContent side="top" className="max-w-56 text-center">
                      A manual reminder. Press ↻ whenever you hear from the artist. It turns amber when it&apos;s been a while.
                    </TooltipContent>
                  </Tooltip>
                </span>
              </TableHead>
            ) : (
              <TableHead className="text-xs font-medium text-muted-foreground">Artpieces</TableHead>
            )}
            <TableHead><span className="sr-only">Actions</span></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {sorted.map(commission => (
            <TableRow key={commission.id}>
              {/* Title */}
              <TableCell className="text-sm font-medium">
                {commission.title ?? <span className="text-muted-foreground">—</span>}
              </TableCell>

              {/* Artist */}
              <TableCell className="text-sm">
                <span className="inline-flex items-center gap-1">
                  {commission.artist_name ?? <span className="text-muted-foreground">—</span>}
                  {commission.artist_link && (
                    <a
                      href={commission.artist_link}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-muted-foreground hover:text-foreground"
                      aria-label="Open artist link"
                    >
                      <ExternalLink className="size-3" />
                    </a>
                  )}
                </span>
              </TableCell>

              {/* Status */}
              <TableCell className="text-sm">
                <span className="inline-flex items-center gap-1">
                  <StatusChip
                    value={commission.status}
                    onChange={(s) => patchStatus(commission.id, s, commission.status, commission.title ?? null, commission.paid_date)}
                    hasError={hasCellError(commission.id, 'status')}
                  />
                  <CellError show={hasCellError(commission.id, 'status')} />
                </span>
              </TableCell>

              {/* Paid */}
              <TableCell className="text-sm">
                <span className={[
                  'inline-flex items-center gap-1 rounded px-1.5 py-1 transition-colors hover:bg-muted',
                  hasCellError(commission.id, 'paid') ? 'ring-1 ring-destructive' : '',
                ].join(' ')}>
                  <Checkbox
                    checked={commission.paid}
                    onCheckedChange={(checked) => patchPaid(commission.id, !!checked)}
                    aria-label="Paid"
                    className="size-5 rounded border-2 border-input data-checked:border-primary data-checked:bg-primary data-checked:text-primary-foreground"
                  />
                  <CellError show={hasCellError(commission.id, 'paid')} />
                </span>
              </TableCell>

              {/* Paid Date */}
              <TableCell className="text-sm">
                <span className="inline-flex items-center gap-1">
                  <Popover>
                    <PopoverTrigger className={[
                      'inline-flex items-center gap-1.5 rounded px-1.5 py-1 text-sm transition-colors hover:bg-muted',
                      hasCellError(commission.id, 'paid_date') ? 'ring-1 ring-destructive' : '',
                    ].join(' ')}>
                      <CalendarIcon className="size-3.5 text-muted-foreground shrink-0" />
                      {commission.paid_date
                        ? new Date(commission.paid_date + 'T12:00:00').toLocaleDateString()
                        : <span className="text-muted-foreground">—</span>
                      }
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                      <Calendar
                        mode="single"
                        selected={commission.paid_date ? new Date(commission.paid_date + 'T12:00:00') : undefined}
                        onSelect={(date) => patchPaidDate(commission.id, date)}
                      />
                    </PopoverContent>
                  </Popover>
                  <CellError show={hasCellError(commission.id, 'paid_date')} />
                </span>
              </TableCell>

              {/* Last Contact (active) / Artpieces (done) */}
              {view === 'active' ? (
                <TableCell className="text-sm">
                  <LastContactCell
                    commission={commission}
                    onStamp={() => stampLastContacted(commission.id)}
                    hasError={hasCellError(commission.id, 'last_contacted_at')}
                  />
                </TableCell>
              ) : (
                <TableCell>
                  <ArtpieceThumbnailStack artpieces={commission.artpieces} />
                </TableCell>
              )}

              {/* Actions */}
              <TableCell>
                <span className="inline-flex items-center gap-1">
                  <Button
                    variant="ghost"
                    size="icon-xs"
                    type="button"
                    onClick={() => onEdit(commission)}
                    aria-label="Edit commission"
                  >
                    <Pencil className="size-3.5" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon-xs"
                    type="button"
                    onClick={() => onDelete(commission)}
                    aria-label="Delete commission"
                    className="text-muted-foreground hover:text-destructive"
                  >
                    <Trash2 className="size-3.5" />
                  </Button>
                </span>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TooltipProvider>
  )
}
