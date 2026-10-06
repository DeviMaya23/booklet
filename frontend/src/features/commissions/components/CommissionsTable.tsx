import { ExternalLink, RefreshCw, ChevronUp, ChevronDown, AlertCircle, CalendarIcon, Pencil, Trash2 } from 'lucide-react'
import StatusChip from './StatusChip'
import {
  Table, TableHeader, TableBody, TableHead, TableRow, TableCell,
} from '@/components/ui/table'
import { Checkbox } from '@/components/ui/checkbox'
import { Popover, PopoverTrigger, PopoverContent } from '@/components/ui/popover'
import { Calendar } from '@/components/ui/calendar'
import { Button } from '@/components/ui/button'
import { type Commission } from '../api/useCommissions'
import { usePatchCommission } from '../api/usePatchCommission'
import { formatRelativeTime } from '../utils/relativeTime'
import { timeTakenDays } from '../utils/timeTakenDays'
import { useCommissionsSort, type SortKey, type SortDir } from '../hooks/useCommissionsSort'
import { useCellErrors } from '../hooks/useCellErrors'
import { useCommissionActions } from '../hooks/useCommissionActions'


interface CommissionsTableProps {
  commissions: Commission[]
  view: 'active' | 'done'
  onEdit: (commission: Commission) => void
  onDelete: (commission: Commission) => void
}

function SortIcon({ column, sortKey, sortDir }: { column: SortKey; sortKey: SortKey; sortDir: SortDir }) {
  if (column !== sortKey) return <ChevronUp className="size-3 opacity-30" />
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
        className="inline-flex items-center gap-1 text-sm font-semibold text-foreground hover:text-foreground/70"
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

export default function CommissionsTable({ commissions, view, onEdit, onDelete }: CommissionsTableProps) {
  const { sorted, sortKey, sortDir, toggleSort } = useCommissionsSort(commissions)
  const { setCellError, hasCellError } = useCellErrors()
  const patchCommission = usePatchCommission()
  const { patchStatus, patchPaid, patchPaidDate, stampLastContacted } = useCommissionActions({
    mutateAsync: patchCommission.mutateAsync,
    setCellError,
  })

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <SortableHead label="Title" column="title" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
          <SortableHead label="Artist" column="artist_name" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
          <SortableHead label="Status" column="status" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
          <SortableHead label="Paid" column="paid" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
          <SortableHead label="Paid Date" column="paid_date" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
          {view === 'active'
            ? <SortableHead label="Last Contact" column="last_contacted_at" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
            : <SortableHead label="Time Taken" column="time_taken" sortKey={sortKey} sortDir={sortDir} onToggle={toggleSort} />
          }
          <TableHead><span className="sr-only">Actions</span></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody className="text-xs">
        {sorted.map(commission => {
          const days = timeTakenDays(commission)
          return (
            <TableRow key={commission.id}>
              {/* Title */}
              <TableCell className="font-medium">
                {commission.title ?? <span className="text-muted-foreground">—</span>}
              </TableCell>

              {/* Artist */}
              <TableCell>
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
              <TableCell>
                <span className="inline-flex items-center gap-1">
                  <StatusChip
                    value={commission.status}
                    onChange={(s) => patchStatus(commission.id, s)}
                    hasError={hasCellError(commission.id, 'status')}
                  />
                  <CellError show={hasCellError(commission.id, 'status')} />
                </span>
              </TableCell>

              {/* Paid */}
              <TableCell>
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
              <TableCell>
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

              {/* Last Contact (active) / Time Taken (done) */}
              {view === 'active' ? (
                <TableCell>
                  <span className={[
                    'inline-flex items-center gap-1.5 rounded px-1.5 py-1 transition-colors',
                    hasCellError(commission.id, 'last_contacted_at') ? 'ring-1 ring-destructive' : '',
                  ].join(' ')}>
                    <span className="text-sm">
                      {commission.last_contacted_at
                        ? formatRelativeTime(new Date(commission.last_contacted_at))
                        : <span className="text-muted-foreground">—</span>
                      }
                    </span>
                    <Button
                      variant="ghost"
                      size="icon-xs"
                      type="button"
                      onClick={() => stampLastContacted(commission.id)}
                      aria-label="Mark contacted now"
                    >
                      <RefreshCw />
                    </Button>
                    <CellError show={hasCellError(commission.id, 'last_contacted_at')} />
                  </span>
                </TableCell>
              ) : (
                <TableCell>
                  <span className="text-sm">
                    {days !== null ? `${days} days` : <span className="text-muted-foreground">—</span>}
                  </span>
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
          )
        })}
      </TableBody>
    </Table>
  )
}
