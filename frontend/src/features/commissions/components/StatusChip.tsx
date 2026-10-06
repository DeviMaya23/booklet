import { ChevronDown } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from '@/components/ui/dropdown-menu'
import { STATUS_OPTIONS, STATUS_LABELS, type CommissionStatus } from '../statusOptions'

interface StatusChipProps {
  value: string
  onChange: (status: CommissionStatus) => void
  hasError?: boolean
  disabled?: boolean
  triggerClassName?: string
}

export default function StatusChip({ value, onChange, hasError, disabled, triggerClassName }: StatusChipProps) {
  const defaultTriggerClassName = [
    'inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors',
    'border hover:bg-muted',
    hasError ? 'border-destructive' : 'border-border',
    disabled ? 'cursor-not-allowed opacity-50' : '',
  ].join(' ')

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        disabled={disabled}
        className={triggerClassName ?? defaultTriggerClassName}
      >
        {STATUS_LABELS[value as CommissionStatus] ?? value}
        <ChevronDown className="size-3 opacity-60" />
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        {STATUS_OPTIONS.map((s) => (
          <DropdownMenuItem key={s} onClick={() => onChange(s)}>
            {STATUS_LABELS[s]}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
