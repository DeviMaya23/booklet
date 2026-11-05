import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import DeleteCharacterDialog from './DeleteCharacterDialog'

function setup(props: Partial<Parameters<typeof DeleteCharacterDialog>[0]> = {}) {
  const onCancel = vi.fn()
  const onConfirm = vi.fn()
  render(<DeleteCharacterDialog open pending={false} onCancel={onCancel} onConfirm={onConfirm} {...props} />)
  return { onCancel, onConfirm }
}

describe('DeleteCharacterDialog', () => {
  it('asks to confirm and confirms', async () => {
    const { onConfirm } = setup()

    expect(screen.getByText('Delete character?')).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Delete' }))

    expect(onConfirm).toHaveBeenCalledTimes(1)
  })

  it('cancels from the Cancel button', async () => {
    const { onCancel, onConfirm } = setup()

    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    expect(onCancel).toHaveBeenCalledTimes(1)
    expect(onConfirm).not.toHaveBeenCalled()
  })

  it('disables Delete while the request is pending', () => {
    setup({ pending: true })

    expect(screen.getByRole('button', { name: 'Delete' })).toBeDisabled()
  })

  it('renders nothing when closed', () => {
    setup({ open: false })

    expect(screen.queryByText('Delete character?')).not.toBeInTheDocument()
  })
})
