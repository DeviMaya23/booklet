import { Star, Pencil } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { type Artist, type ArtistLink } from '../api/useArtists'
import { getHostname } from '../lib/linkUtils'

interface ArtistsListProps {
  artists: Artist[]
  onEditClick: (artist: Artist) => void
}

function LinkChip({ link }: { link: ArtistLink }) {
  const label = getHostname(link.url)
  const tooltipSecondLine = link.is_primary ? 'Main link · opens in a new tab' : 'Opens in a new tab'

  return (
    <Tooltip>
      <TooltipTrigger
        render={
          <a
            href={link.url}
            target="_blank"
            rel="noreferrer"
            onClick={(e) => e.stopPropagation()}
            className="inline-flex items-center gap-1 rounded-md bg-muted px-2 py-0.5 text-xs font-medium text-foreground"
          >
            {link.is_primary && <Star className="size-3 fill-current" />}
            {label}
          </a>
        }
      />
      <TooltipContent side="bottom">
        <div className="flex flex-col gap-0.5">
          <span>{link.url}</span>
          <span className="text-background/70">{tooltipSecondLine}</span>
        </div>
      </TooltipContent>
    </Tooltip>
  )
}

function LinkChips({ links }: { links: ArtistLink[] }) {
  if (links.length === 0) {
    return <span className="text-xs text-muted-foreground">No links</span>
  }

  const sorted = [...links].sort((a, b) => (b.is_primary ? 1 : 0) - (a.is_primary ? 1 : 0))
  const visible = sorted.slice(0, 3)
  const overflow = sorted.length - 3

  return (
    <TooltipProvider>
      <div className="flex flex-wrap items-center gap-1">
        {visible.map((link) => (
          <LinkChip key={link.id} link={link} />
        ))}
        {overflow > 0 && (
          <span className="inline-flex items-center rounded-md bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground">
            +{overflow}
          </span>
        )}
      </div>
    </TooltipProvider>
  )
}

export default function ArtistsList({ artists, onEditClick }: ArtistsListProps) {
  if (artists.length === 0) return null

  return (
    <Table className="table-fixed">
      <TableHeader>
        <TableRow>
          <TableHead className="w-[15%]">Name</TableHead>
          <TableHead className="w-[40%]">Links</TableHead>
          <TableHead className="w-[40%]">Notes</TableHead>
          <TableHead className="w-px" />
        </TableRow>
      </TableHeader>
      <TableBody>
        {artists.map((artist) => (
          <TableRow key={artist.id}>
            <TableCell className="font-semibold">{artist.name}</TableCell>
            <TableCell>
              <LinkChips links={artist.links} />
            </TableCell>
            <TableCell className="truncate text-sm text-muted-foreground">
              {artist.notes ?? '—'}
            </TableCell>
            <TableCell className="w-px">
              <Button
                variant="ghost"
                size="icon-sm"
                onClick={() => onEditClick(artist)}
                title="Edit artist"
              >
                <Pencil className="size-4" />
                <span className="sr-only">Edit artist</span>
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
