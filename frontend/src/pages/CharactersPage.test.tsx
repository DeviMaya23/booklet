import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Routes, Route, useParams } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import CharactersPage from './CharactersPage'
import { createWrapper } from '@/test/queryClient'
import { mockApiRoutes } from '@/test/mockApi'

vi.mock('@kinde-oss/kinde-auth-react', () => ({
  useKindeAuth: vi.fn().mockReturnValue({ getToken: vi.fn().mockResolvedValue('tok') }),
}))
vi.mock('@/lib/api', () => ({ apiFetch: vi.fn() }))
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

const tanelle = {
  id: 'c-1',
  name: 'Tanelle',
  avatar_url: null,
  notes: null,
  is_public: false,
  verse_id: null,
  folders: [],
  created_at: '',
  updated_at: '',
}

function DetailStub() {
  const { id } = useParams()
  return <div>Detail page for {id}</div>
}

function renderPage() {
  const Wrapper = createWrapper()
  return render(
    <Wrapper>
      <MemoryRouter initialEntries={['/app/characters']}>
        <Routes>
          <Route path="/app/characters" element={<CharactersPage />} />
          <Route path="/app/characters/:id" element={<DetailStub />} />
        </Routes>
      </MemoryRouter>
    </Wrapper>,
  )
}

beforeEach(() => {
  vi.clearAllMocks()
  mockApiRoutes({ 'GET /characters': { body: [tanelle] } })
})

describe('CharactersPage', () => {
  it('navigates to the detail page when a card is clicked, without opening a modal', async () => {
    renderPage()

    const caption = await screen.findByText('Tanelle')
    await userEvent.click(caption.closest('figure')!.firstElementChild!)

    expect(await screen.findByText('Detail page for c-1')).toBeInTheDocument()
    expect(screen.queryByText('New character')).not.toBeInTheDocument()
  })

  it('does not navigate when the card menu is used', async () => {
    renderPage()
    await screen.findByText('Tanelle')

    const trigger = screen.getByRole('button', { name: '' })
    await userEvent.click(trigger)

    expect(await screen.findByText('Delete')).toBeInTheDocument()
    expect(screen.queryByText(/Detail page for/)).not.toBeInTheDocument()
  })

  it('opens the empty create modal from "New Character"', async () => {
    renderPage()
    await screen.findByText('Tanelle')

    await userEvent.click(screen.getByRole('button', { name: /new character/i }))

    expect(await screen.findByText('New character')).toBeInTheDocument()
    expect(screen.getByLabelText(/name/i)).toHaveValue('')
  })
})
