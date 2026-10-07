import { useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { toast } from 'sonner'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useCommissions, type Commission } from '@/features/commissions/api/useCommissions'
import { useDeleteCommission } from '@/features/commissions/api/useDeleteCommission'
import { usePatchCommission } from '@/features/commissions/api/usePatchCommission'
import CommissionsTable from '@/features/commissions/components/CommissionsTable'
import CommissionFormModal from '@/features/commissions/components/CommissionFormModal'
import DeleteCommissionDialog from '@/features/commissions/components/DeleteCommissionDialog'

type Tab = 'active' | 'finished'

export default function CommissionsPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const tab: Tab = searchParams.get('tab') === 'finished' ? 'finished' : 'active'

  const { data: commissions, isLoading, isError } = useCommissions()
  const patchCommission = usePatchCommission()

  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<Commission | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Commission | null>(null)

  const deleteCommission = useDeleteCommission()

  async function handleDeleteConfirm() {
    if (!deleteTarget) return
    try {
      await deleteCommission.mutateAsync(deleteTarget.id)
      setDeleteTarget(null)
    } catch {
      toast.error('Failed to delete commission')
    }
  }

  const allCommissions = commissions ?? []
  const activeCount = allCommissions.filter(c => c.status === 'waitlist' || c.status === 'wip').length
  const finishedCount = allCommissions.filter(c => c.status === 'done').length

  const filtered = allCommissions.filter(c =>
    tab === 'active'
      ? c.status === 'waitlist' || c.status === 'wip'
      : c.status === 'done'
  )

  function handleCommissionDone(id: string, prevStatus: string, title: string | null) {
    toast(`${title ?? 'Commission'} moved to Finished`, {
      action: {
        label: 'Undo',
        onClick: async () => {
          try {
            await patchCommission.mutateAsync({ id, status: prevStatus })
          } catch {
            toast.error('Failed to undo — please edit the commission manually')
          }
        },
      },
    })
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <div role="tablist" aria-label="Commission view" className="flex items-center gap-1 rounded-lg bg-muted p-1">
          <button
            role="tab"
            aria-selected={tab === 'active'}
            type="button"
            onClick={() => setSearchParams({ tab: 'active' })}
            className={[
              'inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
              tab === 'active'
                ? 'bg-background shadow-sm text-foreground'
                : 'text-muted-foreground hover:text-foreground',
            ].join(' ')}
          >
            Active
            <span className={[
              'rounded-full px-1.5 py-0.5 text-xs',
              tab === 'active' ? 'bg-muted text-foreground' : 'bg-muted/60 text-muted-foreground',
            ].join(' ')}>
              {isLoading ? '—' : activeCount}
            </span>
          </button>
          <button
            role="tab"
            aria-selected={tab === 'finished'}
            type="button"
            onClick={() => setSearchParams({ tab: 'finished' })}
            className={[
              'inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
              tab === 'finished'
                ? 'bg-background shadow-sm text-foreground'
                : 'text-muted-foreground hover:text-foreground',
            ].join(' ')}
          >
            Finished
            <span className={[
              'rounded-full px-1.5 py-0.5 text-xs',
              tab === 'finished' ? 'bg-muted text-foreground' : 'bg-muted/60 text-muted-foreground',
            ].join(' ')}>
              {isLoading ? '—' : finishedCount}
            </span>
          </button>
        </div>

        <Button onClick={() => setCreateOpen(true)}>
          <Plus className="size-4" />
          New commission
        </Button>
      </div>

      {isLoading && <p className="text-sm text-muted-foreground">Loading…</p>}
      {isError && <p className="text-sm text-destructive">Failed to load commissions.</p>}

      {!isLoading && !isError && filtered.length === 0 && (
        <p className="text-sm text-muted-foreground">No commissions here yet.</p>
      )}

      {!isLoading && !isError && filtered.length > 0 && (
        <CommissionsTable
          commissions={filtered}
          view={tab === 'active' ? 'active' : 'done'}
          onEdit={(c) => setEditTarget(c)}
          onDelete={(c) => setDeleteTarget(c)}
          onCommissionDone={handleCommissionDone}
        />
      )}

      <CommissionFormModal
        open={createOpen}
        onOpenChange={setCreateOpen}
        mode="create"
      />

      <CommissionFormModal
        open={editTarget !== null}
        onOpenChange={(open) => { if (!open) setEditTarget(null) }}
        mode="edit"
        commission={editTarget ?? undefined}
      />

      {deleteTarget && (
        <DeleteCommissionDialog
          open={deleteTarget !== null}
          onOpenChange={(open) => { if (!open) setDeleteTarget(null) }}
          onConfirm={handleDeleteConfirm}
          isPending={deleteCommission.isPending}
        />
      )}
    </div>
  )
}
