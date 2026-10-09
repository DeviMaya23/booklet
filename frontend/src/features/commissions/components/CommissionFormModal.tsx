import { useEffect, useRef, useState } from 'react'
import { CalendarIcon, ChevronDown, ChevronRight, Plus, X } from 'lucide-react'
import { toast } from 'sonner'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Checkbox } from '@/components/ui/checkbox'
import { Popover, PopoverTrigger, PopoverContent } from '@/components/ui/popover'
import { Calendar } from '@/components/ui/calendar'
import { type Artist } from '@/features/artists/api/useArtists'
import ArtistCombobox from '@/features/artists/components/ArtistCombobox'
import ArtistFormModal from '@/features/artists/components/ArtistFormModal'
import ArtistCharacterFilter from '@/components/ArtistCharacterFilter'
import { useArtpieces } from '@/features/artpieces/api/useArtpieces'
import { type Commission } from '../api/useCommissions'
import { useCommission, type CommissionArtpieceSummary } from '../api/useCommission'
import { useCreateCommission } from '../api/useCreateCommission'
import { useUpdateCommission } from '../api/useUpdateCommission'
import { useReplaceArtpieces } from '../api/useReplaceArtpieces'
import StatusChip from './StatusChip'
import { type CommissionStatus } from '../statusOptions'

function formatDateForApi(date: Date | undefined): string | null {
  if (!date) return null
  return [
    date.getFullYear(),
    String(date.getMonth() + 1).padStart(2, '0'),
    String(date.getDate()).padStart(2, '0'),
  ].join('-')
}

function parseDateFromApi(str: string | null): Date | undefined {
  if (!str) return undefined
  return new Date(str + 'T12:00:00')
}

interface CommissionFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  mode: 'create' | 'edit'
  commission?: Commission
}

export default function CommissionFormModal({ open, onOpenChange, mode, commission }: CommissionFormModalProps) {
  const [title, setTitle] = useState('')
  const [selectedArtist, setSelectedArtist] = useState<Artist | null>(null)
  const [notes, setNotes] = useState('')
  const [status, setStatus] = useState<CommissionStatus>('waitlist')
  const [priceInput, setPriceInput] = useState('')
  const [paid, setPaid] = useState(false)
  const [paidDate, setPaidDate] = useState<Date | undefined>()
  const [finishDate, setFinishDate] = useState<Date | undefined>()

  const [artpiecesOpen, setArtpiecesOpen] = useState(false)
  const [stripItems, setStripItems] = useState<CommissionArtpieceSummary[]>([])
  const [search, setSearch] = useState('')
  const [filterArtist, setFilterArtist] = useState<Artist | null>(null)
  const [filterCharacters, setFilterCharacters] = useState<{ id: string; name: string }[]>([])

  const [artistModalOpen, setArtistModalOpen] = useState(false)

  const initialArtpieceIds = useRef<string[]>([])

  const artpiecesQuery = useArtpieces()
  const createCommission = useCreateCommission()
  const updateCommission = useUpdateCommission()
  const replaceArtpieces = useReplaceArtpieces()

  const isSubmitting = createCommission.isPending || updateCommission.isPending || replaceArtpieces.isPending

  const commissionId = mode === 'edit' && commission ? commission.id : null
  const { data: commissionDetail, isLoading: detailLoading } = useCommission(commissionId)

  function populateFromCommission(c: Commission) {
    setTitle(c.title ?? '')
    setNotes(c.notes ?? '')
    setStatus((c.status as CommissionStatus) ?? 'waitlist')
    setPriceInput(c.price != null ? String(c.price) : '')
    setPaid(c.paid)
    setPaidDate(parseDateFromApi(c.paid_date))
    setFinishDate(parseDateFromApi(c.finish_date))
  }

  function resetForm() {
    setTitle('')
    setSelectedArtist(null)
    setNotes('')
    setStatus('waitlist')
    setPriceInput('')
    setPaid(false)
    setPaidDate(undefined)
    setFinishDate(undefined)
    setArtpiecesOpen(false)
    setStripItems([])
    setSearch('')
    setFilterArtist(null)
    setFilterCharacters([])
    initialArtpieceIds.current = []
  }

  useEffect(() => {
    if (!open) {
      resetForm()
      return
    }
    if (mode === 'edit' && commission) {
      populateFromCommission(commission)
      // Artist pre-populate: Commission has artist_id + artist_name — build a minimal Artist object
      if (commission.artist_id && commission.artist_name) {
        setSelectedArtist({
          id: commission.artist_id,
          name: commission.artist_name,
          notes: null,
          links: [],
          created_at: '',
          updated_at: '',
        })
      } else {
        setSelectedArtist(null)
      }
    }
    // create mode: defaults already set
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  // When commissionDetail loads in edit mode, populate the artpiece strip
  useEffect(() => {
    if (mode === 'edit' && commissionDetail && open) {
      setStripItems(commissionDetail.artpieces)
      initialArtpieceIds.current = commissionDetail.artpieces.map((a) => a.id)
    }
  }, [commissionDetail, mode, open])

  const currentIds = stripItems.map((a) => a.id)

  function artpieceSetChanged(): boolean {
    const initial = new Set(initialArtpieceIds.current)
    const current = new Set(currentIds)
    if (initial.size !== current.size) return true
    for (const id of current) {
      if (!initial.has(id)) return true
    }
    return false
  }

  const availableArtpieces = (artpiecesQuery.data ?? []).filter((a) => {
    // Exclude already-in-strip
    if (currentIds.includes(a.id)) return false
    // Exclude attached to a different commission
    if (a.commission_id !== null && a.commission_id !== commissionId) return false
    // Title search
    if (search && !a.title?.toLowerCase().includes(search.toLowerCase())) return false
    // Artist filter
    if (filterArtist && a.artist_id !== filterArtist.id) return false
    // Character filter (any match)
    if (filterCharacters.length > 0) {
      const artpieceCharIds = a.characters.map((c) => c.id)
      if (!filterCharacters.some((fc) => artpieceCharIds.includes(fc.id))) return false
    }
    return true
  })

  function addToStrip(artpiece: { id: string; thumbnail_url: string | null }) {
    setStripItems((prev) => {
      if (prev.some((a) => a.id === artpiece.id)) return prev
      return [...prev, { id: artpiece.id, thumbnail_url: artpiece.thumbnail_url }]
    })
  }

  function removeFromStrip(id: string) {
    setStripItems((prev) => prev.filter((a) => a.id !== id))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    try {
      const price = priceInput.trim() ? parseFloat(priceInput) : null
      const artpieceIds = currentIds

      if (mode === 'create') {
        await createCommission.mutateAsync({
          title: title.trim() || null,
          artistId: selectedArtist?.id ?? null,
          status,
          price: isNaN(price as number) ? null : price,
          paid,
          paidDate: formatDateForApi(paidDate),
          finishDate: formatDateForApi(finishDate),
          notes: notes.trim() || null,
          artpieceIds,
        })
        toast.success('Commission created')
      } else if (mode === 'edit' && commission) {
        await updateCommission.mutateAsync({
          id: commission.id,
          title: title.trim() || null,
          artistId: selectedArtist?.id ?? null,
          status,
          price: isNaN(price as number) ? null : price,
          paid,
          paidDate: formatDateForApi(paidDate),
          finishDate: formatDateForApi(finishDate),
          notes: notes.trim() || null,
        })
        if (artpieceSetChanged()) {
          await replaceArtpieces.mutateAsync({ id: commission.id, artpieceIds })
        }
        toast.success('Commission updated')
      }

      onOpenChange(false)
    } catch {
      toast.error(mode === 'create' ? 'Failed to create commission' : 'Failed to update commission')
    }
  }

  const isLoading = mode === 'edit' && detailLoading && open

  return (
    <>
      <Dialog open={open} onOpenChange={(next) => { if (!isSubmitting) onOpenChange(next) }}>
        <DialogContent className="flex max-h-[calc(100dvh-2rem)] max-w-md flex-col overflow-hidden">
          <DialogHeader>
            <DialogTitle>{mode === 'create' ? 'New commission' : 'Edit commission'}</DialogTitle>
          </DialogHeader>

          <form id="commission-form" onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto pr-1">
            {/* Title */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="commission-title">Title</label>
              <Input
                id="commission-title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Commission title"
                disabled={isSubmitting}
              />
            </div>

            {/* Artist */}
            <label className="flex flex-col gap-1.5">
              <span className="text-sm font-medium">Artist</span>
              <div className="flex gap-2">
                <div className="flex-1">
                  <ArtistCombobox value={selectedArtist} onChange={setSelectedArtist} disabled={isSubmitting} />
                </div>
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  onClick={() => setArtistModalOpen(true)}
                  disabled={isSubmitting}
                >
                  <Plus className="size-4" />
                </Button>
              </div>
            </label>

            {/* Notes */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium" htmlFor="commission-notes">Notes</label>
              <Textarea
                id="commission-notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Notes about this commission…"
                rows={3}
                disabled={isSubmitting}
              />
            </div>

            <hr className="border-border" />

            {/* Status / Price / Paid / Paid Date / Finish Date */}
            <div className="flex flex-col gap-4">
              {/* Row 1: Status | Price */}
              <div className="grid grid-cols-2 gap-3">
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium">Status</label>
                  <StatusChip
                    value={status}
                    onChange={setStatus}
                    disabled={isSubmitting}
                    triggerClassName="flex w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50"
                  />
                </div>
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium">Price</label>
                  <Input
                    type="number"
                    min="0"
                    step="0.01"
                    value={priceInput}
                    onChange={(e) => setPriceInput(e.target.value)}
                    placeholder="0.00"
                    disabled={isSubmitting}
                  />
                </div>
              </div>

              {/* Row 2: Paid | Paid Date */}
              <div className="grid grid-cols-2 gap-3">
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium">Paid</label>
                  <button
                    type="button"
                    onClick={() => !isSubmitting && setPaid((v) => !v)}
                    className="flex items-center gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50"
                    disabled={isSubmitting}
                  >
                    <Checkbox
                      checked={paid}
                      className="size-4 rounded border-input data-checked:border-primary data-checked:bg-primary data-checked:text-primary-foreground pointer-events-none"
                    />
                    <span className="text-muted-foreground">Mark as paid</span>
                  </button>
                </div>
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium">Paid date</label>
                  <Popover>
                    <PopoverTrigger
                      disabled={isSubmitting}
                      className="flex w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      {paidDate
                        ? <span>{paidDate.toLocaleDateString()}</span>
                        : <span className="text-muted-foreground">Pick date</span>
                      }
                      <CalendarIcon className="size-4 text-muted-foreground shrink-0" />
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                      <Calendar mode="single" selected={paidDate} onSelect={setPaidDate} />
                    </PopoverContent>
                  </Popover>
                </div>
              </div>

              {/* Row 3: Finish Date (full width) */}
              <div className="flex flex-col gap-1.5">
                <label className="text-sm font-medium">Finish date</label>
                <Popover>
                  <PopoverTrigger
                    disabled={isSubmitting}
                    className="flex w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-xs hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    {finishDate
                      ? <span>{finishDate.toLocaleDateString()}</span>
                      : <span className="text-muted-foreground">Pick date</span>
                    }
                    <CalendarIcon className="size-4 text-muted-foreground shrink-0" />
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar mode="single" selected={finishDate} onSelect={setFinishDate} />
                  </PopoverContent>
                </Popover>
              </div>
            </div>

            <hr className="border-border" />

            {/* Artpieces section */}
            <div className="flex flex-col gap-2">
              <button
                type="button"
                onClick={() => setArtpiecesOpen((v) => !v)}
                className="flex items-center gap-1.5 text-sm font-medium"
              >
                {artpiecesOpen ? <ChevronDown className="size-4" /> : <ChevronRight className="size-4" />}
                Artpieces
                {currentIds.length > 0 && (
                  <span className="ml-1 rounded-full bg-muted px-1.5 py-0.5 text-xs">{currentIds.length}</span>
                )}
              </button>

              {artpiecesOpen && (
                <div className="flex flex-col gap-3 rounded-md border border-border p-3">
                  {/* Strip */}
                  {isLoading ? (
                    <div className="flex gap-2">
                      {[1, 2, 3].map((i) => (
                        <div key={i} className="size-14 animate-pulse rounded bg-muted" />
                      ))}
                    </div>
                  ) : stripItems.length === 0 ? (
                    <p className="text-xs text-muted-foreground">No artpieces attached yet.</p>
                  ) : (
                    <div className="flex flex-wrap gap-2">
                      {stripItems.map((item) => (
                        <div key={item.id} className="relative">
                          <div className="size-14 overflow-hidden rounded border border-border bg-muted">
                            {item.thumbnail_url ? (
                              <img src={item.thumbnail_url} alt="" className="size-full object-cover" />
                            ) : (
                              <div className="size-full bg-muted" />
                            )}
                          </div>
                          <button
                            type="button"
                            onClick={() => removeFromStrip(item.id)}
                            className="absolute -right-1.5 -top-1.5 flex size-4 items-center justify-center rounded-full bg-background border border-border text-muted-foreground hover:text-foreground"
                            aria-label="Remove artpiece"
                          >
                            <X className="size-2.5" />
                          </button>
                        </div>
                      ))}
                    </div>
                  )}

                  {/* Search + Filter */}
                  <div className="flex gap-2">
                    <Input
                      placeholder="Search by title…"
                      value={search}
                      onChange={(e) => setSearch(e.target.value)}
                    />
                    <ArtistCharacterFilter
                      artist={filterArtist}
                      onArtistChange={setFilterArtist}
                      characters={filterCharacters}
                      onCharactersChange={setFilterCharacters}
                    />
                  </div>

                  {/* Search results */}
                  <div className="flex flex-col gap-1 max-h-40 overflow-y-auto">
                      {artpiecesQuery.isPending ? (
                        <p className="text-xs text-muted-foreground">Loading…</p>
                      ) : artpiecesQuery.isError ? (
                        <p className="text-xs text-destructive">Failed to load artpieces.</p>
                      ) : availableArtpieces.length === 0 ? (
                        <p className="text-xs text-muted-foreground">No artpieces found.</p>
                      ) : (
                        availableArtpieces.map((a) => (
                          <button
                            key={a.id}
                            type="button"
                            onClick={() => addToStrip(a)}
                            className="flex items-center gap-2 rounded px-2 py-1.5 text-left text-sm hover:bg-muted"
                          >
                            <div className="size-8 shrink-0 overflow-hidden rounded border border-border bg-muted">
                              {a.thumbnail_url ? (
                                <img src={a.thumbnail_url} alt="" className="size-full object-cover" />
                              ) : (
                                <div className="size-full bg-muted" />
                              )}
                            </div>
                            <span className="truncate">{a.title ?? <span className="text-muted-foreground">Untitled</span>}</span>
                          </button>
                        ))
                      )}
                  </div>
                </div>
              )}
            </div>
          </form>

          <DialogFooter>
            <Button
              variant="outline"
              type="button"
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button type="submit" form="commission-form" disabled={isSubmitting}>
              {isSubmitting ? 'Saving…' : 'Save'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ArtistFormModal
        open={artistModalOpen}
        onOpenChange={setArtistModalOpen}
        onCreated={(artist) => setSelectedArtist(artist)}
      />
    </>
  )
}
