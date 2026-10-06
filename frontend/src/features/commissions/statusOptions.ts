export const STATUS_OPTIONS = ['waitlist', 'wip', 'done'] as const
export type CommissionStatus = typeof STATUS_OPTIONS[number]
export const STATUS_LABELS: Record<CommissionStatus, string> = {
  waitlist: 'Waitlist',
  wip: 'WIP',
  done: 'Done',
}
