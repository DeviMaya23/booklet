import { Loader2 } from 'lucide-react'
import { getFileCategory } from '../lib/fileCategory'

interface FileThumbProps {
  thumbnailUrl: string | null
  extension: string | null
  thumbnailGenState?: 'pending' | 'done' | 'failed' | 'not_applicable'
}

export default function FileThumb({ thumbnailUrl, extension, thumbnailGenState }: FileThumbProps) {
  if (thumbnailUrl) {
    return <img src={thumbnailUrl} alt="" className="h-full w-full object-cover" draggable={false} />
  }

  if (thumbnailGenState === 'pending') {
    return (
      <div className="flex h-full w-full items-center justify-center bg-muted">
        <Loader2 className="h-6 w-6 text-muted-foreground animate-spin" />
      </div>
    )
  }

  const { icon: Icon, label, tinted, cssVar } = getFileCategory(extension)

  const containerStyle = tinted && cssVar
    ? { backgroundColor: `var(${cssVar}-bg)`, color: `var(${cssVar}-fg)` }
    : undefined

  return (
    <div
      className={`flex h-full w-full flex-col items-center justify-center gap-1 ${tinted ? '' : 'bg-muted text-muted-foreground'}`}
      style={containerStyle}
    >
      <div className="w-[28%] aspect-square flex items-center justify-center">
        <Icon className="h-full w-full" strokeWidth={1.5} />
      </div>
      {label && (
        <span className="max-w-[80%] truncate rounded px-1 py-0.5 text-[10px] font-medium leading-none opacity-80">
          {label}
        </span>
      )}
    </div>
  )
}
