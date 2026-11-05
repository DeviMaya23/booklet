import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import DiscardChangesDialog from './DiscardChangesDialog'

function setup(open = true) {
  const onKeepEditing = vi.fn()
  const onDiscard = vi.fn()
  render(<DiscardChangesDialog open={open} onKeepEditing={onKeepEditing} onDiscard={onDiscard} />)
  return { onKeepEditing, onDiscard }
}

describe('DiscardChangesDialog', () => {
  it('discards from the Discard button without also keeping the edits', async () => {
    const { onDiscard, onKeepEditing } = setup()

    expect(screen.getByText('Discard changes?')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Discard' }))

    expect(onDiscard).toHaveBeenCalledTimes(1)
    expect(onKeepEditing).not.toHaveBeenCalled()
  })

  it('keeps editing from the Keep editing button', async () => {
    const { onKeepEditing, onDiscard } = setup()

    await userEvent.click(screen.getByRole('button', { name: 'Keep editing' }))

    expect(onKeepEditing).toHaveBeenCalledTimes(1)
    expect(onDiscard).not.toHaveBeenCalled()
  })

  it('keeps editing when dismissed with Escape', async () => {
    const { onKeepEditing } = setup()

    await userEvent.keyboard('{Escape}')

    expect(onKeepEditing).toHaveBeenCalledTimes(1)
  })

  it('renders nothing when closed', () => {
    setup(false)

    expect(screen.queryByText('Discard changes?')).not.toBeInTheDocument()
  })
})
