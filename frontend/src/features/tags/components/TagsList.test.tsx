import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import TagsList from './TagsList'

const noop = vi.fn()

function renderList(props: Partial<React.ComponentProps<typeof TagsList>>) {
  render(
    <MemoryRouter>
      <TagsList
        tags={[]}
        isEmpty={false}
        editingId={null}
        renameError={null}
        isSaving={false}
        onEdit={noop}
        onCancelEdit={noop}
        onSubmitName={noop}
        onDelete={noop}
        {...props}
      />
    </MemoryRouter>,
  )
}

describe('TagsList', () => {
  it('renders one row per tag', () => {
    renderList({
      tags: [
        { id: 'a', name: 'outfit', artpiece_count: 22 },
        { id: 'b', name: 'sketch', artpiece_count: 14 },
      ],
    })

    expect(screen.getAllByRole('listitem')).toHaveLength(2)
    expect(screen.getByText('22 artpieces →')).toBeInTheDocument()
    expect(screen.getByText('14 artpieces →')).toBeInTheDocument()
  })

  it('shows the empty state when the user has no tags', () => {
    renderList({ isEmpty: true })

    expect(screen.getByText('No tags yet')).toBeInTheDocument()
    expect(screen.getByText(/edit page will show up here/)).toBeInTheDocument()
    expect(screen.queryByRole('listitem')).not.toBeInTheDocument()
  })

  it('shows no empty state and no rows when search matches nothing', () => {
    renderList({ tags: [], isEmpty: false })

    expect(screen.queryByText('No tags yet')).not.toBeInTheDocument()
    expect(screen.queryByRole('listitem')).not.toBeInTheDocument()
  })

  it('puts only the editing row into edit mode', () => {
    renderList({
      tags: [
        { id: 'a', name: 'outfit', artpiece_count: 1 },
        { id: 'b', name: 'sketch', artpiece_count: 1 },
      ],
      editingId: 'b',
    })

    expect(screen.getAllByRole('textbox')).toHaveLength(1)
    expect(screen.getByRole('textbox')).toHaveValue('sketch')
  })
})
