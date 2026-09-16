import { useRef, useState } from 'react'
import { X } from 'lucide-react'

interface TokenInputProps<T extends { id: string; name: string }> {
  items: T[]
  onChange: (items: T[]) => void
  disabled?: boolean
  suggestions?: T[]
  placeholder?: string
  // When provided, enables free-text entry: comma key and blur commit a raw
  // string via this function. Returning null means "don't add" (e.g. empty
  // after trim, or duplicate name).
  createFromText?: (raw: string) => T | null
}

export default function TokenInput<T extends { id: string; name: string }>({
  items,
  onChange,
  disabled,
  suggestions = [],
  placeholder,
  createFromText,
}: TokenInputProps<T>) {
  const [val, setVal] = useState('')
  const [selectedIndex, setSelectedIndex] = useState(-1)
  const [showDropdown, setShowDropdown] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)
  const blurTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const filtered = suggestions.filter(
    (s) =>
      val.trim().length > 0 &&
      s.name.toLowerCase().includes(val.trim().toLowerCase()) &&
      !items.some((i) => i.id === s.id),
  )

  const dropdownVisible = showDropdown && filtered.length > 0

  const add = (item: T) => {
    onChange([...items, item])
    setVal('')
    setSelectedIndex(-1)
    setShowDropdown(false)
  }

  const commitRaw = (raw: string) => {
    if (!createFromText) return
    const item = createFromText(raw)
    if (!item) {
      setVal('')
      return
    }
    add(item)
  }

  const remove = (target: T) => onChange(items.filter((i) => i !== target))

  const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      if (!dropdownVisible) return
      setSelectedIndex((i) => (i + 1) % filtered.length)
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      if (!dropdownVisible) return
      setSelectedIndex((i) => (i <= 0 ? filtered.length - 1 : i - 1))
    } else if (e.key === 'Enter') {
      e.preventDefault()
      if (dropdownVisible && selectedIndex >= 0) {
        add(filtered[selectedIndex])
      } else if (createFromText) {
        commitRaw(val)
      }
    } else if (e.key === ',' && createFromText) {
      e.preventDefault()
      commitRaw(val)
    } else if (e.key === 'Escape') {
      setVal('')
      setSelectedIndex(-1)
      setShowDropdown(false)
    } else if (e.key === 'Backspace' && val === '' && items.length) {
      remove(items[items.length - 1])
    }
  }

  const onValChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setVal(e.target.value)
    setSelectedIndex(-1)
    setShowDropdown(true)
  }

  const onBlur = () => {
    blurTimerRef.current = setTimeout(() => {
      if (createFromText && val.trim()) {
        commitRaw(val)
      } else {
        setVal('')
      }
      setShowDropdown(false)
      setSelectedIndex(-1)
    }, 150)
  }

  const onSuggestionMouseDown = (item: T) => {
    // Cancel the blur timer so the dropdown doesn't close before click fires
    if (blurTimerRef.current) clearTimeout(blurTimerRef.current)
    add(item)
    inputRef.current?.focus()
  }

  return (
    <div className="relative">
      <div
        onClick={() => !disabled && inputRef.current?.focus()}
        className={`flex flex-wrap gap-1.5 px-3 py-1.5 border border-input rounded-md bg-transparent min-h-9 items-center cursor-text shadow-xs focus-within:border-ring focus-within:ring-[3px] focus-within:ring-ring/50 ${disabled ? 'opacity-50 cursor-default' : ''}`}
      >
        {items.map((item) => (
          <span
            key={item.id || item.name}
            className="inline-flex items-center gap-1 bg-muted rounded-md px-2 py-0.5 text-sm select-none"
          >
            {item.name}
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation()
                remove(item)
              }}
              className="text-muted-foreground hover:text-foreground transition-colors ml-0.5"
              aria-label={`Remove ${item.name}`}
            >
              <X className="w-2.5 h-2.5" />
            </button>
          </span>
        ))}
        <input
          ref={inputRef}
          value={val}
          onChange={onValChange}
          onKeyDown={onKeyDown}
          onFocus={() => setShowDropdown(true)}
          onBlur={onBlur}
          placeholder={items.length === 0 ? placeholder : ''}
          disabled={disabled}
          className="border-none outline-none bg-transparent text-sm min-w-[70px] flex-1 py-0.5 px-0.5"
        />
      </div>

      {dropdownVisible && (
        <ul className="absolute z-10 left-0 right-0 mt-1 bg-popover border border-border rounded-md shadow-md overflow-hidden">
          {filtered.map((s, i) => (
            <li
              key={s.id}
              onMouseDown={() => onSuggestionMouseDown(s)}
              className={`px-3 py-1.5 text-sm cursor-pointer select-none ${
                i === selectedIndex
                  ? 'bg-accent text-accent-foreground'
                  : 'hover:bg-muted'
              }`}
            >
              {s.name}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
