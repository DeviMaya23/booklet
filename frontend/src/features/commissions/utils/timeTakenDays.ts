import { type Commission } from '../api/useCommissions'

export function timeTakenDays(commission: Commission): number | null {
  if (!commission.finish_date || !commission.paid_date) return null
  const finish = new Date(commission.finish_date).getTime()
  const paid = new Date(commission.paid_date).getTime()
  const days = Math.round((finish - paid) / (1000 * 60 * 60 * 24))
  return days < 0 ? null : days
}
