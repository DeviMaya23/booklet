import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { ChooseFoldersPanel } from './ChooseFoldersPanel'
import type { PublicFolder } from '../api/usePublicFolders'

const folderA: PublicFolder = { id: 'fa', name: 'Plant', token: 'tok-plant' }
const folderB: PublicFolder = { id: 'fb', name: 'Cottagecore', token: 'tok-cottage' }
const folderC: PublicFolder = { id: 'fc', name: 'Tavern', token: 'tok-tavern' }

const noop = () => {}

function renderPanel(overrides: Partial<Parameters<typeof ChooseFoldersPanel>[0]> = {}) {
  const defaults = {
    selectedIds: new Set<string>(),
    onAdd: noop,
    onRemove: noop,
    onClose: noop,
    folders: [folderA, folderB, folderC],
    isLoading: false,
    isError: false,
    onRefresh: noop,
  }
  return render(<ChooseFoldersPanel {...defaults} {...overrides} />)
}

describe('ChooseFoldersPanel', () => {
  it('renders all folders with unchecked checkboxes when none are selected', () => {
    renderPanel()
    expect(screen.getByRole('checkbox', { name: 'Plant' })).not.toBeChecked()
    expect(screen.getByRole('checkbox', { name: 'Cottagecore' })).not.toBeChecked()
    expect(screen.getByRole('checkbox', { name: 'Tavern' })).not.toBeChecked()
  })

  it('renders selected folders with checked checkboxes', () => {
    renderPanel({ selectedIds: new Set(['fa', 'fc']) })
    expect(screen.getByRole('checkbox', { name: 'Plant' })).toBeChecked()
    expect(screen.getByRole('checkbox', { name: 'Cottagecore' })).not.toBeChecked()
    expect(screen.getByRole('checkbox', { name: 'Tavern' })).toBeChecked()
  })

  it('calls onAdd when unchecked folder is toggled', async () => {
    const onAdd = vi.fn()
    renderPanel({ onAdd })
    await userEvent.click(screen.getByRole('checkbox', { name: 'Plant' }))
    expect(onAdd).toHaveBeenCalledWith(folderA)
  })

  it('calls onRemove when checked folder is toggled', async () => {
    const onRemove = vi.fn()
    renderPanel({ selectedIds: new Set(['fa']), onRemove })
    await userEvent.click(screen.getByRole('checkbox', { name: 'Plant' }))
    expect(onRemove).toHaveBeenCalledWith('fa')
  })

  it('filters folder list case-insensitively as user types', async () => {
    renderPanel()
    await userEvent.type(screen.getByPlaceholderText(/search folders/i), 'plant')
    expect(screen.getByText('Plant')).toBeInTheDocument()
    expect(screen.queryByText('Cottagecore')).not.toBeInTheDocument()
    expect(screen.queryByText('Tavern')).not.toBeInTheDocument()
  })

  it('shows all folders again when search is cleared', async () => {
    renderPanel()
    const input = screen.getByPlaceholderText(/search folders/i)
    await userEvent.type(input, 'plant')
    await userEvent.clear(input)
    expect(screen.getByText('Plant')).toBeInTheDocument()
    expect(screen.getByText('Cottagecore')).toBeInTheDocument()
    expect(screen.getByText('Tavern')).toBeInTheDocument()
  })

  it('shows empty state when folder list is empty', () => {
    renderPanel({ folders: [] })
    expect(screen.getByText('No public folders yet')).toBeInTheDocument()
    const link = screen.getByText('Open Bookleaf').closest('a')
    expect(link).toHaveAttribute('target', '_blank')
  })

  it('shows error state when isError is true', () => {
    renderPanel({ isError: true, folders: null })
    expect(screen.getByText(/couldn't reach bookleaf/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /retry/i })).toBeInTheDocument()
  })

  it('calls onRefresh when Retry is clicked in error state', async () => {
    const onRefresh = vi.fn()
    renderPanel({ isError: true, folders: null, onRefresh })
    await userEvent.click(screen.getByRole('button', { name: /retry/i }))
    expect(onRefresh).toHaveBeenCalled()
  })

  it('calls onClose when X button is clicked', async () => {
    const onClose = vi.fn()
    renderPanel({ onClose })
    await userEvent.click(screen.getByRole('button', { name: /close folder panel/i }))
    expect(onClose).toHaveBeenCalled()
  })
})
