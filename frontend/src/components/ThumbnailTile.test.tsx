import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { vi, describe, it, expect } from 'vitest'
import ThumbnailTile, { type ThumbnailTileMenuItem } from './ThumbnailTile'

function renderTile(props: Partial<Parameters<typeof ThumbnailTile>[0]> = {}) {
  return render(
    <MemoryRouter>
      <ThumbnailTile
        imageUrl="https://cdn.example.com/thumb.jpg"
        label="Chev"
        href="/app/artpieces?artist=123"
        {...props}
      />
    </MemoryRouter>,
  )
}

describe('ThumbnailTile', () => {
  it('renders image when URL is provided', () => {
    renderTile()
    expect(screen.getByRole('img', { name: 'Chev' })).toHaveAttribute(
      'src',
      'https://cdn.example.com/thumb.jpg',
    )
  })

  it('renders placeholder when imageUrl is null', () => {
    renderTile({ imageUrl: null })
    expect(screen.queryByRole('img', { name: 'Chev' })).not.toBeInTheDocument()
  })

  it('renders label below the image', () => {
    renderTile()
    expect(screen.getByText('Chev')).toBeInTheDocument()
  })

  it('renders sublabel when provided', () => {
    renderTile({ sublabel: 'chel' })
    expect(screen.getByText('chel')).toBeInTheDocument()
  })

  it('does not render sublabel when absent', () => {
    renderTile()
    expect(screen.queryByText('chel')).not.toBeInTheDocument()
  })

  it('does not render action button when menuItems is absent', () => {
    renderTile()
    expect(screen.queryByRole('button', { name: 'More options' })).not.toBeInTheDocument()
  })

  it('renders action button when menuItems are provided', () => {
    const menuItems: ThumbnailTileMenuItem[] = [
      { label: 'Delete', variant: 'destructive', onClick: vi.fn() },
    ]
    renderTile({ menuItems })
    expect(screen.getByRole('button', { name: 'More options' })).toBeInTheDocument()
  })

})
