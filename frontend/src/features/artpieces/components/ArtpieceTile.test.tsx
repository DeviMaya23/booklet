import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect } from 'vitest'
import ArtpieceTile from './ArtpieceTile'

const noop = () => {}

const defaultProps = {
  title: 'Sunrise',
  artistName: 'Alice',
  imageUrl: 'https://cdn.example.com/cover.jpg',
  showDetails: true,
  onClick: noop,
  onInfoClick: noop,
}

describe('ArtpieceTile', () => {
  it('renders thumbnail when imageUrl is present', () => {
    render(<ArtpieceTile {...defaultProps} />)
    const img = screen.getByRole('img', { name: 'Sunrise' })
    expect(img).toBeInTheDocument()
    expect(img).toHaveAttribute('src', 'https://cdn.example.com/cover.jpg')
  })

  it('renders placeholder when imageUrl is null', () => {
    render(<ArtpieceTile {...defaultProps} imageUrl={null} />)
    expect(screen.queryByRole('img', { name: 'Sunrise' })).not.toBeInTheDocument()
    expect(screen.getByRole('img', { name: 'No cover image' })).toBeInTheDocument()
  })

  it('renders caption row with title and artist when showDetails is true', () => {
    render(<ArtpieceTile {...defaultProps} showDetails={true} />)
    expect(screen.getByText('Sunrise')).toBeInTheDocument()
    expect(screen.getByText('Alice')).toBeInTheDocument()
  })

  it('hides caption row when showDetails is false', () => {
    render(<ArtpieceTile {...defaultProps} showDetails={false} />)
    expect(screen.queryByText('Sunrise')).not.toBeInTheDocument()
    expect(screen.queryByText('Alice')).not.toBeInTheDocument()
  })

  it('calls onClick when the tile is clicked', async () => {
    const onClick = vi.fn()
    render(<ArtpieceTile {...defaultProps} onClick={onClick} />)
    await userEvent.click(screen.getByRole('button', { name: 'View Sunrise' }))
    expect(onClick).toHaveBeenCalledTimes(1)
  })

  it('calls onInfoClick when ⓘ button is clicked and does not fire onClick', async () => {
    const onClick = vi.fn()
    const onInfoClick = vi.fn()
    render(<ArtpieceTile {...defaultProps} onClick={onClick} onInfoClick={onInfoClick} />)
    await userEvent.click(screen.getByRole('button', { name: 'Open details' }))
    expect(onInfoClick).toHaveBeenCalledTimes(1)
    expect(onClick).not.toHaveBeenCalled()
  })
})
