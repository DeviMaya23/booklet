import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { FolderList } from './FolderList'
import type { PublicFolder } from '../api/usePublicFolders'

const folder1: PublicFolder = { id: 'f1', name: 'Plant', token: 'tok-plant' }
const folder2: PublicFolder = { id: 'f2', name: 'Cottagecore', token: 'tok-cottage' }
const invalidFolder = { id: 'f3', name: 'Old refs' }

const noop = () => {}

describe('FolderList', () => {
  it('renders valid folder rows with name and share link', () => {
    render(
      <FolderList
        selectedFolders={[folder1, folder2]}
        invalidFolders={[]}
        onRemove={noop}
        onDismissInvalid={noop}
        onOpenPanel={noop}
      />,
    )
    expect(screen.getByText('Plant')).toBeInTheDocument()
    expect(screen.getByText('Cottagecore')).toBeInTheDocument()
    const link = screen.getByRole('link', { name: /open plant in bookleaf/i })
    expect(link).toHaveAttribute('href', expect.stringContaining('/share/tok-plant'))
    expect(link).toHaveAttribute('target', '_blank')
  })

  it('renders warning row for invalid folder with subtitle and Remove button', () => {
    render(
      <FolderList
        selectedFolders={[]}
        invalidFolders={[invalidFolder]}
        onRemove={noop}
        onDismissInvalid={noop}
        onOpenPanel={noop}
      />,
    )
    expect(screen.getByText('Old refs')).toBeInTheDocument()
    expect(screen.getByText('No longer public in Bookleaf.')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /remove/i })).toBeInTheDocument()
  })

  it('calls onRemove with the folder id when remove button is clicked', async () => {
    const onRemove = vi.fn()
    render(
      <FolderList
        selectedFolders={[folder1]}
        invalidFolders={[]}
        onRemove={onRemove}
        onDismissInvalid={noop}
        onOpenPanel={noop}
      />,
    )
    await userEvent.click(screen.getByRole('button', { name: /remove plant/i }))
    expect(onRemove).toHaveBeenCalledWith('f1')
  })

  it('calls onDismissInvalid with the folder id when Remove is clicked on a warning row', async () => {
    const onDismissInvalid = vi.fn()
    render(
      <FolderList
        selectedFolders={[]}
        invalidFolders={[invalidFolder]}
        onRemove={noop}
        onDismissInvalid={onDismissInvalid}
        onOpenPanel={noop}
      />,
    )
    await userEvent.click(screen.getByRole('button', { name: /remove/i }))
    expect(onDismissInvalid).toHaveBeenCalledWith('f3')
  })

  it('calls onOpenPanel when Choose folders button is clicked', async () => {
    const onOpenPanel = vi.fn()
    render(
      <FolderList
        selectedFolders={[]}
        invalidFolders={[]}
        onRemove={noop}
        onDismissInvalid={noop}
        onOpenPanel={onOpenPanel}
      />,
    )
    await userEvent.click(screen.getByRole('button', { name: /choose folders/i }))
    expect(onOpenPanel).toHaveBeenCalled()
  })
})
