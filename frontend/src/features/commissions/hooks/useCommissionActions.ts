import { type PatchCommissionParams } from '../api/usePatchCommission'

interface Deps {
  mutateAsync: (params: PatchCommissionParams) => Promise<unknown>
  setCellError: (id: string, field: string, error: boolean) => void
}

export function useCommissionActions({ mutateAsync, setCellError }: Deps) {
  async function patchStatus(id: string, status: string) {
    try {
      await mutateAsync({ id, status })
      setCellError(id, 'status', false)
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
