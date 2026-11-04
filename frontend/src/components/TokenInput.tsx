import { useEffect, useRef, useState } from 'react'
import { X } from 'lucide-react'

// Max suggestions listed when the dropdown opens with empty text. Typed
// matches are never capped.
const FOCUS_LIST_LIMIT = 10

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
  // Opt-in: list suggestions as soon as the input is focused, before typing.
  openOnFocus?: boolean
  // Opt-in async creation. Enter or comma with no suggestion highlighted, or
  // choosing the create row, calls this with the trimmed text. A resolved item
  // is added as a chip; null keeps the text in place. Free text is never
  // created on blur. Ignored while a previous call is still pending.
  onCreate?: (raw: string) => Promise<T | null>
  // Whether the create row is offered for the current trimmed text. Defaults
  // to true for any non-empty text when onCreate is provided.
  canCreate?: (raw: string) => boolean
}

export default function TokenInput<T extends { id: string; name: string }>({
  items,
  onChange,
  disabled,
  suggestions = [],
  placeholder,
  createFromText,
  openOnFocus,
  onCreate,
  canCreate,
}: TokenInputProps<T>) {
  const [val, setVal] = useState('')
  const [selectedIndex, setSelectedIndex] = useState(-1)
  const [showDropdown, setShowDropdown] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)
  const [pending, setPending] = useState(false)
  const pendingRef = useRef(false)
  const blurTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  // A pending blur timer must not outlive the component.
  useEffect(
    () => () => {
      if (blurTimerRef.current) clearTimeout(blurTimerRef.current)
    },
    [],
  )

  const trimmed = val.trim()
  const matches = suggestions.filter(
    (s) =>
      (trimmed.length > 0 ? s.name.toLowerCase().includes(trimmed.toLowerCase()) : !!openOnFocus) &&
      !items.some((i) => i.id === s.id),
  )
  const truncated = trimmed.length === 0 && matches.length > FOCUS_LIST_LIMIT
  const filtered = truncated ? matches.slice(0, FOCUS_LIST_LIMIT) : matches

  const showCreateRow = !!onCreate && trimmed.length > 0 && (canCreate ? canCreate(trimmed) : true)
  const rowCount = filtered.length + (showCreateRow ? 1 : 0)
  const dropdownVisible = showDropdown && rowCount > 0

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

  const runCreate = async () => {
    if (!onCreate || !trimmed || pendingRef.current) return
    pendingRef.current = true
    setPending(true)
    try {
      const item = await onCreate(trimmed)
      if (!item) return
      if (items.some((i) => i.id === item.id)) {
        setVal('')
        setSelectedIndex(-1)
        setShowDropdown(false)
      } else {
        add(item)
      }
    } finally {
      pendingRef.current = false
      setPending(false)
    }
  }

  // Chips are locked while an async create is pending: it resolves against the
  // items of the render that started it, so a removal now would be undone.
  const remove = (target: T) => {
    if (pendingRef.current) return
    onChange(items.filter((i) => i !== target))
  }

  const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      if (!dropdownVisible) return
      setSelectedIndex((i) => (i + 1) % rowCount)
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      if (!dropdownVisible) return
      setSelectedIndex((i) => (i <= 0 ? rowCount - 1 : i - 1))
    } else if (e.key === 'Enter') {
      e.preventDefault()
      if (dropdownVisible && selectedIndex >= 0 && selectedIndex < filtered.length) {
        add(filtered[selectedIndex])
      } else if (onCreate) {
        void runCreate()
      } else if (createFromText) {
        commitRaw(val)
      }
    } else if (e.key === ',' && onCreate) {
      e.preventDefault()
      void runCreate()
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
      // While a create is pending the typed text stays (it may still be
      // needed if the create fails), but the dropdown still closes.
      if (!pendingRef.current) {
        if (createFromText && val.trim()) {
          commitRaw(val)
        } else {
          setVal('')
        }
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

  const onCreateRowMouseDown = () => {
    if (blurTimerRef.current) clearTimeout(blurTimerRef.current)
    void runCreate()
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
          onFocus={() => {
            // Refocusing before the blur timer fires must cancel it, or it
            // would clear the text and close the list the user is using.
            if (blurTimerRef.current) clearTimeout(blurTimerRef.current)
            setShowDropdown(true)
          }}
          onBlur={onBlur}
          placeholder={items.length === 0 ? placeholder : ''}
          disabled={disabled}
          readOnly={pending}
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
          {truncated && (
            <li
              // Keep focus on the input so the blur timer doesn't close the list.
              onMouseDown={(e) => e.preventDefault()}
              className="px-3 py-1.5 text-xs text-muted-foreground select-none cursor-default"
            >
              Type to search more…
            </li>
          )}
          {showCreateRow && (
            <li
              onMouseDown={onCreateRowMouseDown}
              className={`px-3 py-1.5 text-sm cursor-pointer select-none ${
                selectedIndex === filtered.length
                  ? 'bg-accent text-accent-foreground'
                  : 'hover:bg-muted'
              }`}
            >
              Create &ldquo;{trimmed}&rdquo;
            </li>
          )}
        </ul>
      )}
    </div>
  )
}
