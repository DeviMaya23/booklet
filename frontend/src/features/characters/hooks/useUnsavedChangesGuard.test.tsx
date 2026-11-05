import { act, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import { useUnsavedChangesGuard } from './useUnsavedChangesGuard'

function Harness({ isDirty, onLeave }: { isDirty: boolean; onLeave: () => void }) {
  const guard = useUnsavedChangesGuard({ isDirty, onLeave })
  return (
    <div>
      <p>{guard.promptOpen ? 'prompt open' : 'prompt closed'}</p>
      <button onClick={guard.requestLeave}>request leave</button>
      <button onClick={guard.discard}>discard</button>
      <button onClick={guard.keepEditing}>keep editing</button>
      <button onClick={guard.allowLeaving}>allow leaving</button>
    </div>
  )
}

function setup(isDirty: boolean) {
  const onLeave = vi.fn()
  const router = createMemoryRouter(
    [
      { path: '/edit', element: <Harness isDirty={isDirty} onLeave={onLeave} /> },
      { path: '/elsewhere', element: <p>elsewhere</p> },
    ],
    { initialEntries: ['/edit'] },
  )
  render(<RouterProvider router={router} />)
  return { router, onLeave }
}

describe('useUnsavedChangesGuard', () => {
  it('leaves straight away when nothing is dirty', async () => {
    const { onLeave } = setup(false)

    await userEvent.click(screen.getByText('request leave'))

    expect(onLeave).toHaveBeenCalledTimes(1)
    expect(screen.getByText('prompt closed')).toBeInTheDocument()
  })

  it('asks first when dirty, and Discard then leaves', async () => {
    const { onLeave } = setup(true)

    await userEvent.click(screen.getByText('request leave'))
    expect(screen.getByText('prompt open')).toBeInTheDocument()
    expect(onLeave).not.toHaveBeenCalled()

    await userEvent.click(screen.getByText('discard'))

    expect(onLeave).toHaveBeenCalledTimes(1)
    expect(screen.getByText('prompt closed')).toBeInTheDocument()
  })

  it('keeps editing: closes the prompt without leaving', async () => {
    const { onLeave } = setup(true)

    await userEvent.click(screen.getByText('request leave'))
    await userEvent.click(screen.getByText('keep editing'))

    expect(screen.getByText('prompt closed')).toBeInTheDocument()
    expect(onLeave).not.toHaveBeenCalled()
  })

  it('blocks in-app navigation while dirty; Discard lets it through', async () => {
    const { router, onLeave } = setup(true)

    await act(async () => {
      void router.navigate('/elsewhere')
    })
    expect(screen.getByText('prompt open')).toBeInTheDocument()
    expect(screen.queryByText('elsewhere')).not.toBeInTheDocument()

    await userEvent.click(screen.getByText('discard'))

    expect(await screen.findByText('elsewhere')).toBeInTheDocument()
    expect(onLeave).not.toHaveBeenCalled()
  })

  it('keep editing resets a blocked navigation', async () => {
    const { router } = setup(true)

    await act(async () => {
      void router.navigate('/elsewhere')
    })
    await userEvent.click(screen.getByText('keep editing'))

    expect(screen.getByText('prompt closed')).toBeInTheDocument()
    expect(router.state.location.pathname).toBe('/edit')
  })

  it('does not block navigation when nothing is dirty', async () => {
    const { router } = setup(false)

    await act(async () => {
      void router.navigate('/elsewhere')
    })

    expect(await screen.findByText('elsewhere')).toBeInTheDocument()
  })

  it('lets a navigation through after allowLeaving, even when dirty', async () => {
    const { router } = setup(true)

    await userEvent.click(screen.getByText('allow leaving'))
    await act(async () => {
      void router.navigate('/elsewhere')
    })

    expect(await screen.findByText('elsewhere')).toBeInTheDocument()
  })
})
