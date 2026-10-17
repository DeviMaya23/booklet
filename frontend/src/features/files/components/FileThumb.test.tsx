import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import FileThumb from './FileThumb'

describe('FileThumb — thumbnail present', () => {
  it('renders the thumbnail image when thumbnailUrl is set', () => {
    const { container } = render(
      <FileThumb thumbnailUrl="https://cdn.example.com/thumb.jpg" extension="jpg" />,
    )
    const img = container.querySelector('img')
    expect(img).not.toBeNull()
    expect(img).toHaveAttribute('src', 'https://cdn.example.com/thumb.jpg')
  })
})

describe('FileThumb — fallback', () => {
  it('renders the category fallback when thumbnailUrl is null', () => {
    const { container } = render(<FileThumb thumbnailUrl={null} extension="mp4" />)
    expect(container.querySelector('img')).toBeNull()
    expect(screen.getByText('MP4')).toBeInTheDocument()
  })

  it('switches to the fallback icon when the image fires onError', () => {
    const { container } = render(
      <FileThumb thumbnailUrl="https://cdn.example.com/thumb.jpg" extension="png" />,
    )
    const img = container.querySelector('img')!
    fireEvent.error(img)
    expect(container.querySelector('img')).toBeNull()
    expect(screen.getByText('PNG')).toBeInTheDocument()
  })
})

describe('FileThumb — imgFailed reset', () => {
  it('retries the image when thumbnailUrl changes after a load failure', () => {
    const { container, rerender } = render(
      <FileThumb thumbnailUrl="https://cdn.example.com/expired.jpg" extension="jpg" />,
    )
    fireEvent.error(container.querySelector('img')!)
    expect(container.querySelector('img')).toBeNull()

    rerender(<FileThumb thumbnailUrl="https://cdn.example.com/fresh.jpg" extension="jpg" />)
    const img = container.querySelector('img')
    expect(img).not.toBeNull()
    expect(img).toHaveAttribute('src', 'https://cdn.example.com/fresh.jpg')
  })

  it('stays on fallback when thumbnailUrl becomes null after a load failure', () => {
    const { container, rerender } = render(
      <FileThumb thumbnailUrl="https://cdn.example.com/expired.jpg" extension="jpg" />,
    )
    fireEvent.error(container.querySelector('img')!)
    rerender(<FileThumb thumbnailUrl={null} extension="jpg" />)
    expect(container.querySelector('img')).toBeNull()
  })
})
