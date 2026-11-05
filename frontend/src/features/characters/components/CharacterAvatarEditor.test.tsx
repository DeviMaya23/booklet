import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import CharacterAvatarEditor from './CharacterAvatarEditor'
import type { AvatarPicker } from '../hooks/useAvatarPicker'

function picker(overrides: Partial<AvatarPicker> = {}): AvatarPicker {
  return {
    localFile: null,
    displayUrl: 'https://img/a.jpg',
    avatarCleared: false,
    canClear: true,
    isChanged: false,
    cropSourceUrl: null,
    dragActive: false,
    dragHandlers: { onDragOver: vi.fn(), onDragLeave: vi.fn(), onDrop: vi.fn() },
    handleFileChange: vi.fn(),
    handlePaste: vi.fn(),
    confirmCrop: vi.fn(),
    cancelCrop: vi.fn(),
    clear: vi.fn(),
    reset: vi.fn(),
    ...overrides,
  }
}

describe('CharacterAvatarEditor', () => {
  it('shows the avatar with a Change overlay and a remove control', () => {
    render(<CharacterAvatarEditor avatar={picker()} />)

    expect(screen.getByAltText('Avatar')).toHaveAttribute('src', 'https://img/a.jpg')
    expect(screen.getByText('Change')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Remove avatar' })).toBeInTheDocument()
  })

  it('shows a prompt and no remove control when there is no avatar', () => {
    render(<CharacterAvatarEditor avatar={picker({ displayUrl: null, canClear: false })} />)

    expect(screen.getByText('Click, drop or paste an image')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Remove avatar' })).not.toBeInTheDocument()
  })

  it('clears the avatar through the picker', async () => {
    const avatar = picker()
    render(<CharacterAvatarEditor avatar={avatar} />)

    await userEvent.click(screen.getByRole('button', { name: 'Remove avatar' }))

    expect(avatar.clear).toHaveBeenCalledTimes(1)
  })

  it('hands a chosen file to the picker', async () => {
    const avatar = picker()
    render(<CharacterAvatarEditor avatar={avatar} />)
    const input = document.querySelector('input[type="file"]') as HTMLInputElement

    await userEvent.upload(input, new File(['png'], 'a.png', { type: 'image/png' }))

    expect(avatar.handleFileChange).toHaveBeenCalledTimes(1)
  })

  it('disables the tile and the remove control while busy', () => {
    render(<CharacterAvatarEditor avatar={picker()} disabled />)

    expect(screen.getByRole('button', { name: 'Change avatar' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'Remove avatar' })).toBeDisabled()
  })
})
