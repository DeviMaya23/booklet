import { useState } from 'react'
import { ArrowLeftRight, Plus } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { useCommissions, type Commission } from '@/features/commissions/api/useCommissions'
import { useDeleteCommission } from '@/features/commissions/api/useDeleteCommission'
import CommissionsTable from '@/features/commissions/components/CommissionsTable'
import CommissionFormModal from '@/features/commissions/components/CommissionFormModal'
import DeleteCommissionDialog from '@/features/commissions/components/DeleteCommissionDialog'

type View = 'active' | 'done'

export default function CommissionsPage() {
  const [view, setView] = useState<View>('active')
  const { data: commissions, isLoading, isError } = useCommissions()

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

  const filtered = (commissions ?? []).filter(c =>
    view === 'active'
      ? c.status === 'waitlist' || c.status === 'wip'
      : c.status === 'done'
  )

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <button
          type="button"
          onClick={() => setView(v => v === 'active' ? 'done' : 'active')}
          className="inline-flex items-center gap-2 cursor-pointer select-none hover:text-foreground"
        >
          <span className="text-base font-semibold">
            {view === 'active' ? 'Waitlist / In Progress' : 'Finished'}
          </span>
          <ArrowLeftRight className="size-4 text-muted-foreground" />
        </button>
        <Button onClick={() => setCreateOpen(true)}>
          <Plus className="size-4" />
          New Commission
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
          view={view}
          onEdit={(c) => setEditTarget(c)}
          onDelete={(c) => setDeleteTarget(c)}
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
