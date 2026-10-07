import { type PatchCommissionParams } from '../api/usePatchCommission'

interface Deps {
  mutateAsync: (params: PatchCommissionParams) => Promise<unknown>
  setCellError: (id: string, field: string, error: boolean) => void
  onStatusDone?: (id: string, prevStatus: string, title: string | null) => void
}

export function useCommissionActions({ mutateAsync, setCellError, onStatusDone }: Deps) {
  async function patchStatus(id: string, newStatus: string, prevStatus: string, title: string | null, currentPaidDate: string | null) {
    try {
      const params: PatchCommissionParams = { id, status: newStatus }
      if (newStatus === 'done' && !currentPaidDate) {
        const today = new Date()
        params.paidDate = [
          today.getFullYear(),
          String(today.getMonth() + 1).padStart(2, '0'),
          String(today.getDate()).padStart(2, '0'),
        ].join('-')
      }
      await mutateAsync(params)
      setCellError(id, 'status', false)
      if (newStatus === 'done') onStatusDone?.(id, prevStatus, title)
    } catch {
      setCellError(id, 'status', true)
    }
  }

  async function patchPaid(id: string, paid: boolean) {
    try {
      const params: PatchCommissionParams = paid
        ? { id, paid, lastContactedAt: new Date().toISOString() }
        : { id, paid }
      await mutateAsync(params)
      setCellError(id, 'paid', false)
    } catch {
      setCellError(id, 'paid', true)
    }
  }

  async function patchPaidDate(id: string, date: Date | undefined) {
    const paidDate = date
      ? [
          date.getFullYear(),
          String(date.getMonth() + 1).padStart(2, '0'),
          String(date.getDate()).padStart(2, '0'),
        ].join('-')
      : null
    try {
      await mutateAsync({ id, paidDate })
      setCellError(id, 'paid_date', false)
    } catch {
      setCellError(id, 'paid_date', true)
    }
  }

  async function stampLastContacted(id: string) {
    try {
      await mutateAsync({ id, lastContactedAt: new Date().toISOString() })
      setCellError(id, 'last_contacted_at', false)
    } catch {
      setCellError(id, 'last_contacted_at', true)
    }
  }

  return { patchStatus, patchPaid, patchPaidDate, stampLastContacted }
}
