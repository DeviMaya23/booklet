import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import TagRow from './TagRow'

function renderRow(props: Partial<React.ComponentProps<typeof TagRow>> = {}) {
  const handlers = {
    onEdit: vi.fn(),
    onCancelEdit: vi.fn(),
    onSubmitName: vi.fn(),
    onDelete: vi.fn(),
  }
  render(
    <MemoryRouter>
      <ul>
        <TagRow tag={{ id: 't1', name: 'outfit', artpiece_count: 22 }} editing={false} {...handlers} {...props} />
      </ul>
    </MemoryRouter>,
  )
  return handlers
}

describe('TagRow display', () => {
  it('shows the name and a usage link to the filtered gallery', () => {
    renderRow()

    expect(screen.getByText('outfit')).toBeInTheDocument()
    const link = screen.getByRole('link', { name: /22 artpieces/ })
    expect(link).toHaveAttribute('href', '/app/artpieces?tags=t1')
    expect(screen.queryByText('Unused')).not.toBeInTheDocument()
  })

  it('uses the singular label for one artpiece', () => {
    renderRow({ tag: { id: 't1', name: 'outfit', artpiece_count: 1 } })

    expect(screen.getByRole('link', { name: /1 artpiece →/ })).toBeInTheDocument()
  })

  it('badges an unused tag and does not link its usage label', () => {
    renderRow({ tag: { id: 't1', name: 'outfit', artpiece_count: 0 } })

    expect(screen.getByText('Unused')).toBeInTheDocument()
    expect(screen.getByText('0 artpieces')).toBeInTheDocument()
    expect(screen.queryByRole('link')).not.toBeInTheDocument()
  })

  it('fires edit and delete actions', async () => {
    const h = renderRow()

    await userEvent.click(screen.getByRole('button', { name: 'Edit outfit' }))
    await userEvent.click(screen.getByRole('button', { name: 'Delete outfit' }))

    expect(h.onEdit).toHaveBeenCalledTimes(1)
    expect(h.onDelete).toHaveBeenCalledTimes(1)
  })
})

describe('TagRow editing', () => {
  it('shows a focused, pre-filled input with controls and hint', () => {
    renderRow({ editing: true })

    const input = screen.getByRole('textbox', { name: 'Tag name' })
    expect(input).toHaveValue('outfit')
    expect(input).toHaveFocus()
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument()
    expect(screen.getByText('Enter to save · Esc to cancel')).toBeInTheDocument()
  })

  it('submits the trimmed new name on Enter', async () => {
    const h = renderRow({ editing: true })
    const input = screen.getByRole('textbox', { name: 'Tag name' })

    await userEvent.clear(input)
    await userEvent.type(input, '  outfits {Enter}')

    expect(h.onSubmitName).toHaveBeenCalledWith('outfits')
  })

  it('submits through the save button', async () => {
    const h = renderRow({ editing: true })
    const input = screen.getByRole('textbox', { name: 'Tag name' })

    await userEvent.type(input, 's')
    await userEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(h.onSubmitName).toHaveBeenCalledWith('outfits')
  })

  it('cancels on Escape and on the cancel button without submitting', async () => {
    const h = renderRow({ editing: true })

    await userEvent.type(screen.getByRole('textbox', { name: 'Tag name' }), 'x{Escape}')
    await userEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    expect(h.onCancelEdit).toHaveBeenCalledTimes(2)
    expect(h.onSubmitName).not.toHaveBeenCalled()
  })

  it('leaves edit mode without a request when the name is unchanged', async () => {
    const h = renderRow({ editing: true })

    await userEvent.type(screen.getByRole('textbox', { name: 'Tag name' }), '{Enter}')

    expect(h.onSubmitName).not.toHaveBeenCalled()
    expect(h.onCancelEdit).toHaveBeenCalledTimes(1)
  })

  it('keeps the input open and sends nothing for a blank name', async () => {
    const h = renderRow({ editing: true })
    const input = screen.getByRole('textbox', { name: 'Tag name' })

    await userEvent.clear(input)
    await userEvent.type(input, '   {Enter}')

    expect(h.onSubmitName).not.toHaveBeenCalled()
    expect(h.onCancelEdit).not.toHaveBeenCalled()
    expect(screen.getByRole('textbox', { name: 'Tag name' })).toBeInTheDocument()
  })

  it('shows an error message while editing', () => {
    renderRow({ editing: true, error: 'Failed to rename tag' })

    expect(screen.getByRole('alert')).toHaveTextContent('Failed to rename tag')
  })
})
