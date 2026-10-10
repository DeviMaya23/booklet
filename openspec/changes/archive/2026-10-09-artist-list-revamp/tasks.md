## 1. Rewrite ArtistsList component

- [x] 1.1 Replace the flex-div layout with `<Table>` / `<TableHeader>` / `<TableBody>` / `<TableRow>` / `<TableCell>` from `/components/ui/table`
- [x] 1.2 Add table header row with columns: Name, Links, Notes (edit action column has no header)
- [x] 1.3 Render artist name in bold in the Name cell
- [x] 1.4 Move the edit icon button into the last cell; remove the copy-link button entirely
- [x] 1.5 Add Notes cell: display `notes` text with `max-w-0 w-full truncate` for ellipsis, or "—" when null

## 2. Implement link chips

- [x] 2.1 Sort links primary-first before rendering: `[...artist.links].sort((a, b) => (b.is_primary ? 1 : 0) - (a.is_primary ? 1 : 0))`
- [x] 2.2 Extract hostname from each URL via `new URL(url).hostname`, falling back to the raw URL in a try/catch
- [x] 2.3 Render up to 3 chips; for each chip render an `<a>` anchor with `target="_blank" rel="noreferrer"`
- [x] 2.4 Style the primary chip with a star icon (`lucide-react` `Star`) and outlined border; style other chips in muted filled style
- [x] 2.5 Render a non-interactive `+N` chip when `links.length > 3`, where N = `links.length - 3`
- [x] 2.6 Render "No links" muted placeholder text when the artist has no links

## 3. Add tooltips to chips

- [x] 3.1 Wrap each chip anchor in `<Tooltip>` / `<TooltipTrigger>` / `<TooltipContent>` from `/components/ui/tooltip`
- [x] 3.2 Tooltip content: full URL on first line; "Main link · opens in a new tab" on second line for primary, "Opens in a new tab" for others

## 4. Update tests

- [x] 4.1 Remove the two obsolete tests ("enables the copy link button", "disables the copy link button")
- [x] 4.2 Add test: primary link chip is rendered first regardless of array order from API
- [x] 4.3 Add test: chips display hostname only, not the full URL
- [x] 4.4 Add test: at most 3 chips shown; a `+N` element appears when there are more than 3 links
- [x] 4.5 Add test: "No links" text is shown when artist has no links
- [x] 4.6 Add test: "—" is shown in the Notes cell when notes is null
- [x] 4.7 Keep existing test: renders nothing (empty DOM) when artists array is empty

## 5. Build and lint

- [x] 5.1 Run `npm run build` and fix any type or build errors
- [x] 5.2 Run `npm run lint` and fix any lint warnings or errors
