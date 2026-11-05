import { useRef } from 'react'
import { Image, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { AvatarPicker } from '../hooks/useAvatarPicker'

interface CharacterAvatarEditorProps {
  avatar: AvatarPicker
  disabled?: boolean
}

/** The avatar tile of the edit form: click, drop or paste to replace, X to remove. */
export default function CharacterAvatarEditor({ avatar, disabled }: CharacterAvatarEditorProps) {
  const fileInputRef = useRef<HTMLInputElement>(null)

  return (
    <div className="relative">
      <button
        type="button"
        aria-label="Change avatar"
        className={cn(
          'relative flex aspect-square w-full cursor-pointer items-center justify-center overflow-hidden rounded-2xl bg-muted transition-colors',
          avatar.dragActive && 'ring-2 ring-ring',
        )}
        onClick={() => fileInputRef.current?.click()}
        {...avatar.dragHandlers}
        disabled={disabled}
      >
        {avatar.displayUrl ? (
          <>
            <img src={avatar.displayUrl} alt="Avatar" className="size-full object-cover" />
            <span className="absolute inset-x-2 bottom-2 rounded-lg bg-black/55 py-1 text-center text-sm font-medium text-white">
              Change
            </span>
          </>
        ) : (
          <span className="flex flex-col items-center gap-1 text-sm text-muted-foreground">
            <Image className="size-8" />
            Click, drop or paste an image
          </span>
        )}
      </button>
      {avatar.canClear && (
        <button
          type="button"
          aria-label="Remove avatar"
          className="absolute right-2 top-2 z-10 inline-flex size-7 items-center justify-center rounded-full bg-black/55 text-white hover:bg-black/70"
          onClick={avatar.clear}
          disabled={disabled}
        >
          <X className="size-4" />
        </button>
      )}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png"
        className="hidden"
        onChange={avatar.handleFileChange}
      />
    </div>
  )
}
