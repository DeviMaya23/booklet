import { useState } from 'react'
import { ArrowLeftRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useCommissions } from '@/features/commissions/api/useCommissions'
import CommissionsTable from '@/features/commissions/components/CommissionsTable'

type View = 'active' | 'done'

export default function CommissionsPage() {
  const [view, setView] = useState<View>('active')
  const { data: commissions, isLoading, isError } = useCommissions()

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
        <Button variant="outline" disabled>
          + New Commission
        </Button>
      </div>

      {isLoading && <p className="text-sm text-muted-foreground">Loading…</p>}
      {isError && <p className="text-sm text-destructive">Failed to load commissions.</p>}

      {!isLoading && !isError && filtered.length === 0 && (
        <p className="text-sm text-muted-foreground">No commissions here yet.</p>
      )}

      {!isLoading && !isError && filtered.length > 0 && (
        <CommissionsTable commissions={filtered} view={view} />
      )}
    </div>
  )
}
