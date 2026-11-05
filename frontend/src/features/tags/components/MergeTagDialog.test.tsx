import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import MergeTagDialog from './MergeTagDialog'

const item = { source: { id: 's', name: 'sketches' }, target: { id: 't', name: 'sketch' } }

describe('MergeTagDialog', () => {
  it('explains the consequences of the merge', () => {
    render(<MergeTagDialog item={item} onCancel={vi.fn()} onConfirm={vi.fn()} />)

    expect(screen.getByText('Merge into “sketch”?')).toBeInTheDocument()
    const body = screen.getByRole('alertdialog').textContent ?? ''
    expect(body).toMatch(/Artpieces tagged “sketches” will get “sketch” instead/)
    expect(body).toMatch(/Artpieces that\s+have both keep just “sketch”/)
    expect(body).toMatch(/“sketches” will be deleted/)
    expect(body).toMatch(/This can’t be\s+undone/)
  })

  it('renders nothing when there is no item', () => {
    render(<MergeTagDialog item={null} onCancel={vi.fn()} onConfirm={vi.fn()} />)

    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()
  })

  it('calls onConfirm from "Merge tags" and onCancel from "Cancel"', async () => {
    const onConfirm = vi.fn()
    const onCancel = vi.fn()
    render(<MergeTagDialog item={item} onCancel={onCancel} onConfirm={onConfirm} />)

    await userEvent.click(screen.getByRole('button', { name: 'Merge tags' }))
    expect(onConfirm).toHaveBeenCalledTimes(1)

    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))
    expect(onCancel).toHaveBeenCalledTimes(1)
  })
})
