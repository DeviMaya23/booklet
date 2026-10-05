import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect } from 'vitest'
import ArtpiecesFilterPopover from './ArtpiecesFilterPopover'

vi.mock('@/features/artists/api/useArtists', () => ({
  useArtists: () => ({ data: [], isError: false }),
}))

vi.mock('@/features/characters/api/useCharacters', () => ({
  useCharacters: () => ({ data: [], isError: false }),
}))

const noop = () => {}

function renderPopover() {
  return render(
    <ArtpiecesFilterPopover
      selectedArtist={null}
      onArtistChange={noop}
      selectedCharacters={[]}
      onCharactersChange={noop}
      characterMatch="all"
      onCharacterMatchChange={noop}
    />,
  )
}

describe('ArtpiecesFilterPopover', () => {
  it('shows the filter panel when the trigger is clicked', async () => {
    renderPopover()

    expect(screen.queryByText('Artist')).not.toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: /filter/i }))

    expect(screen.getByText('Artist')).toBeInTheDocument()
    expect(screen.getByText('Characters')).toBeInTheDocument()
  })

  it('closes the filter panel when Escape is pressed', async () => {
    renderPopover()

    await userEvent.click(screen.getByRole('button', { name: /filter/i }))
    expect(screen.getByText('Artist')).toBeInTheDocument()

    await userEvent.keyboard('{Escape}')

    expect(screen.queryByText('Artist')).not.toBeInTheDocument()
  })
})
