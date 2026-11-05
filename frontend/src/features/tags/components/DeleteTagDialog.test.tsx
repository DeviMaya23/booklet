import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import DeleteTagDialog from './DeleteTagDialog'

const tag = (count: number) => ({ id: 't1', name: 'sketch', artpiece_count: count })

describe('DeleteTagDialog', () => {
  it('states how many artpieces will lose the tag', () => {
    render(<DeleteTagDialog item={tag(14)} onCancel={vi.fn()} onConfirm={vi.fn()} />)

    expect(screen.getByText('Delete “sketch”?')).toBeInTheDocument()
    expect(screen.getByText(/14 artpieces use this tag\. They will lose it\./)).toBeInTheDocument()
    expect(screen.getByText(/Nothing else about them changes/)).toBeInTheDocument()
  })

  it('uses singular wording for one artpiece', () => {
    render(<DeleteTagDialog item={tag(1)} onCancel={vi.fn()} onConfirm={vi.fn()} />)

    expect(screen.getByText(/1 artpiece uses this tag\. It will lose it\./)).toBeInTheDocument()
  })

  it('states that no artpieces use an unused tag', () => {
    render(<DeleteTagDialog item={tag(0)} onCancel={vi.fn()} onConfirm={vi.fn()} />)

    expect(screen.getByText(/No artpieces use this tag\./)).toBeInTheDocument()
  })

  it('renders nothing when there is no item', () => {
    render(<DeleteTagDialog item={null} onCancel={vi.fn()} onConfirm={vi.fn()} />)

    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()
  })

  it('calls onConfirm and onCancel from the buttons', async () => {
    const onConfirm = vi.fn()
    const onCancel = vi.fn()
    render(<DeleteTagDialog item={tag(2)} onCancel={onCancel} onConfirm={onConfirm} />)

    await userEvent.click(screen.getByRole('button', { name: 'Delete' }))
    expect(onConfirm).toHaveBeenCalledTimes(1)

    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))
    expect(onCancel).toHaveBeenCalledTimes(1)
  })
})
