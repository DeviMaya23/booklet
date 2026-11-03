import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, beforeEach } from 'vitest'
import { ThemeProvider } from '@/hooks/useTheme'
import AppSection from './AppSection'

function renderSection() {
  return render(
    <ThemeProvider>
      <AppSection />
    </ThemeProvider>,
  )
}

beforeEach(() => {
  localStorage.clear()
  document.documentElement.removeAttribute('data-theme')
})

describe('AppSection', () => {
  it('lists Parchment, Lumen and Sunless with their subtitles', () => {
    renderSection()

    expect(screen.getByRole('radio', { name: 'Parchment theme' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'Lumen theme' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'Sunless theme' })).toBeInTheDocument()
    expect(screen.getByText('Warm parchment')).toBeInTheDocument()
    expect(screen.getByText('Bright and clean')).toBeInTheDocument()
    expect(screen.getByText('Mostly black')).toBeInTheDocument()
  })

  it('preselects the current theme', () => {
    localStorage.setItem('booklet-theme', 'sunless')

    renderSection()

    expect(screen.getByRole('radio', { name: 'Sunless theme' })).toBeChecked()
    expect(screen.getByRole('radio', { name: 'Parchment theme' })).not.toBeChecked()
    expect(screen.getByRole('radio', { name: 'Lumen theme' })).not.toBeChecked()
  })

  it('selecting Lumen selects it and changes the active theme', async () => {
    renderSection()

    await userEvent.click(screen.getByRole('radio', { name: 'Lumen theme' }))

    expect(screen.getByRole('radio', { name: 'Lumen theme' })).toBeChecked()
    expect(screen.getByRole('radio', { name: 'Parchment theme' })).not.toBeChecked()
    expect(document.documentElement.dataset.theme).toBe('lumen')
  })
})
