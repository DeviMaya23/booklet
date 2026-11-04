import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi, describe, it, expect } from 'vitest'
import TokenInput from './TokenInput'

interface Item {
  id: string
  name: string
}

const nature: Item = { id: 'item-nature', name: 'Nature' }
const travel: Item = { id: 'item-travel', name: 'Travel' }

describe('TokenInput — chip add/remove', () => {
  it('selecting a suggestion calls onChange with the item appended', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} suggestions={[nature, travel]} placeholder="Add…" />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'nat')

    const option = await screen.findByText('Nature')
    await userEvent.click(option)

    expect(onChange).toHaveBeenCalledWith([nature])
  })

  it('removing a chip via the remove button calls onChange without that item', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[nature, travel]} onChange={onChange} placeholder="Add…" />)

    await userEvent.click(screen.getByRole('button', { name: /remove nature/i }))

    expect(onChange).toHaveBeenCalledWith([travel])
  })

  it('Backspace on an empty input removes the last item', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[nature, travel]} onChange={onChange} placeholder="Add…" />)

    const input = screen.getByRole('textbox')
    await userEvent.click(input)
    await userEvent.keyboard('{Backspace}')

    expect(onChange).toHaveBeenCalledWith([nature])
  })
})

describe('TokenInput — dropdown keyboard navigation', () => {
  it('pressing ArrowDown then Enter selects the first suggestion', async () => {
    const onChange = vi.fn()
    const suggestions = [nature, { id: 'item-narrative', name: 'Narrative' }]
    render(<TokenInput items={[]} onChange={onChange} suggestions={suggestions} placeholder="Add…" />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'na')
    await userEvent.keyboard('{ArrowDown}{Enter}')

    expect(onChange).toHaveBeenCalledWith([nature])
  })

  it('ArrowDown wraps around to the first option after the last', async () => {
    const onChange = vi.fn()
    const suggestions = [nature, travel]
    render(<TokenInput items={[]} onChange={onChange} suggestions={suggestions} placeholder="Add…" />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'a')
    await userEvent.keyboard('{ArrowDown}{ArrowDown}{ArrowDown}{Enter}')

    expect(onChange).toHaveBeenCalledWith([nature])
  })

  it('already-applied items are not shown in the dropdown', async () => {
    render(<TokenInput items={[nature]} onChange={vi.fn()} suggestions={[nature]} placeholder="Add…" />)

    const input = screen.getByRole('textbox')
    await userEvent.type(input, 'nat')

    expect(screen.queryByRole('listitem')).not.toBeInTheDocument()
  })
})

describe('TokenInput — blur-close behavior', () => {
  it('blurring without selecting a suggestion does not call onChange (no createFromText)', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} suggestions={[nature, travel]} placeholder="Add…" />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'nat')
    await userEvent.tab()

    await waitFor(() => {
      expect(onChange).not.toHaveBeenCalled()
    })
  })
})

describe('TokenInput — createFromText', () => {
  const createFromText = (raw: string): Item | null => {
    const name = raw.trim().toLowerCase().replace(/,/g, '')
    if (!name) return null
    return { id: '', name }
  }

  it('Enter without a selected suggestion commits raw text via createFromText', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} placeholder="Add…" createFromText={createFromText} />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'concept{Enter}')

    expect(onChange).toHaveBeenCalledWith([{ id: '', name: 'concept' }])
  })

  it('comma key commits raw text via createFromText', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} placeholder="Add…" createFromText={createFromText} />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'concept,')

    expect(onChange).toHaveBeenCalledWith([{ id: '', name: 'concept' }])
  })

  it('blur with a non-empty value commits raw text via createFromText', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} placeholder="Add…" createFromText={createFromText} />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.click(input)
    await userEvent.type(input, 'concept')
    await userEvent.tab()

    await waitFor(() => {
      expect(onChange).toHaveBeenCalledWith([{ id: '', name: 'concept' }])
    })
  })

  it('does not call onChange when createFromText returns null for a duplicate', async () => {
    const dedupCreateFromText = (raw: string): Item | null => {
      const name = raw.trim().toLowerCase().replace(/,/g, '')
      if (!name) return null
      if (name === 'nature') return null
      return { id: '', name }
    }
    const onChange = vi.fn()
    render(<TokenInput items={[nature]} onChange={onChange} placeholder="Add…" createFromText={dedupCreateFromText} />)

    const input = screen.getByRole('textbox')
    await userEvent.type(input, 'nature{Enter}')

    expect(onChange).not.toHaveBeenCalled()
  })

  it('does not call onChange when the trimmed value is empty', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} placeholder="Add…" createFromText={createFromText} />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, '   {Enter}')

    expect(onChange).not.toHaveBeenCalled()
  })
})

describe('TokenInput — openOnFocus', () => {
  it('lists suggestions on focus with empty text when enabled', async () => {
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={[nature, travel]} openOnFocus placeholder="Add…" />)

    await userEvent.click(screen.getByPlaceholderText('Add…'))

    expect(await screen.findByText('Nature')).toBeInTheDocument()
    expect(screen.getByText('Travel')).toBeInTheDocument()
  })

  it('does not list suggestions on focus by default', async () => {
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={[nature, travel]} placeholder="Add…" />)

    await userEvent.click(screen.getByPlaceholderText('Add…'))

    expect(screen.queryByText('Nature')).not.toBeInTheDocument()
  })

  it('excludes already-selected items from the focus list', async () => {
    render(<TokenInput items={[nature]} onChange={vi.fn()} suggestions={[nature, travel]} openOnFocus />)

    await userEvent.click(screen.getByRole('textbox'))

    expect(await screen.findByText('Travel')).toBeInTheDocument()
    expect(screen.getAllByText('Nature')).toHaveLength(1) // the chip only
  })
})

describe('TokenInput — async create', () => {
  const created: Item = { id: 'item-new', name: 'wip' }

  it('offers a create row for non-empty text and creates on click', async () => {
    const onChange = vi.fn()
    const onCreate = vi.fn().mockResolvedValue(created)
    render(<TokenInput items={[]} onChange={onChange} suggestions={[nature]} onCreate={onCreate} placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), 'wip')
    await userEvent.click(await screen.findByText(/Create/))

    expect(onCreate).toHaveBeenCalledWith('wip')
    await waitFor(() => expect(onChange).toHaveBeenCalledWith([created]))
  })

  it('hides the create row when canCreate returns false', async () => {
    const onCreate = vi.fn()
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={[nature]} onCreate={onCreate} canCreate={() => false} placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), 'zzz')

    expect(screen.queryByText(/Create/)).not.toBeInTheDocument()
  })

  it('Enter with no suggestion highlighted creates', async () => {
    const onChange = vi.fn()
    const onCreate = vi.fn().mockResolvedValue(created)
    render(<TokenInput items={[]} onChange={onChange} onCreate={onCreate} placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), '  wip  {Enter}')

    expect(onCreate).toHaveBeenCalledWith('wip')
    await waitFor(() => expect(onChange).toHaveBeenCalledWith([created]))
  })

  it('comma creates', async () => {
    const onCreate = vi.fn().mockResolvedValue(created)
    render(<TokenInput items={[]} onChange={vi.fn()} onCreate={onCreate} placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), 'wip,')

    expect(onCreate).toHaveBeenCalledWith('wip')
  })

  it('the create row is reachable with the arrow keys', async () => {
    const onChange = vi.fn()
    const onCreate = vi.fn().mockResolvedValue(created)
    render(<TokenInput items={[]} onChange={onChange} suggestions={[{ id: 'item-wipe', name: 'wipe' }]} onCreate={onCreate} placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), 'wip')
    await userEvent.keyboard('{ArrowDown}{ArrowDown}{Enter}')

    expect(onCreate).toHaveBeenCalledWith('wip')
  })

  it('Enter on a highlighted suggestion selects it instead of creating', async () => {
    const onChange = vi.fn()
    const onCreate = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} suggestions={[nature]} onCreate={onCreate} placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), 'nat')
    await userEvent.keyboard('{ArrowDown}{Enter}')

    expect(onChange).toHaveBeenCalledWith([nature])
    expect(onCreate).not.toHaveBeenCalled()
  })

  it('blank text creates nothing', async () => {
    const onCreate = vi.fn()
    render(<TokenInput items={[]} onChange={vi.fn()} onCreate={onCreate} placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), '   {Enter}')

    expect(onCreate).not.toHaveBeenCalled()
  })

  it('a null result adds nothing and keeps the text', async () => {
    const onChange = vi.fn()
    const onCreate = vi.fn().mockResolvedValue(null)
    render(<TokenInput items={[]} onChange={onChange} onCreate={onCreate} placeholder="Add…" />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'wip{Enter}')

    await waitFor(() => expect(onCreate).toHaveBeenCalled())
    expect(onChange).not.toHaveBeenCalled()
    expect(input).toHaveValue('wip')
  })

  it('a second Enter while creation is pending creates once', async () => {
    let resolve!: (v: Item) => void
    const onCreate = vi.fn().mockReturnValue(new Promise<Item>((r) => { resolve = r }))
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} onCreate={onCreate} placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), 'wip')
    await userEvent.keyboard('{Enter}{Enter}')
    resolve(created)

    await waitFor(() => expect(onChange).toHaveBeenCalledTimes(1))
    expect(onCreate).toHaveBeenCalledTimes(1)
  })

  it('a resolved item that is already selected is not added twice', async () => {
    const onChange = vi.fn()
    const onCreate = vi.fn().mockResolvedValue(nature)
    render(<TokenInput items={[nature]} onChange={onChange} onCreate={onCreate} />)

    const input = screen.getByRole('textbox')
    await userEvent.type(input, 'nature{Enter}')

    await waitFor(() => expect(onCreate).toHaveBeenCalled())
    await waitFor(() => expect(input).toHaveValue(''))
    expect(onChange).not.toHaveBeenCalled()
  })

  it('blur does not create', async () => {
    const onCreate = vi.fn()
    render(<TokenInput items={[]} onChange={vi.fn()} onCreate={onCreate} placeholder="Add…" />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'wip')
    await userEvent.tab()

    await waitFor(() => expect(input).toHaveValue(''))
    expect(onCreate).not.toHaveBeenCalled()
  })

  it('ignores chip removal while a create is pending, so the chip is not resurrected', async () => {
    let resolve!: (v: Item) => void
    const onCreate = vi.fn().mockReturnValue(new Promise<Item>((r) => { resolve = r }))
    const onChange = vi.fn()
    render(<TokenInput items={[nature]} onChange={onChange} onCreate={onCreate} />)

    await userEvent.type(screen.getByRole('textbox'), 'wip{Enter}')
    await userEvent.click(screen.getByRole('button', { name: /remove nature/i }))
    expect(onChange).not.toHaveBeenCalled()

    resolve(created)
    await waitFor(() => expect(onChange).toHaveBeenCalledTimes(1))
    expect(onChange).toHaveBeenCalledWith([nature, created])
  })

  it('allows chip removal again once the create has settled', async () => {
    const onChange = vi.fn()
    const onCreate = vi.fn().mockResolvedValue(null)
    render(<TokenInput items={[nature]} onChange={onChange} onCreate={onCreate} />)

    await userEvent.type(screen.getByRole('textbox'), 'wip{Enter}')
    await waitFor(() => expect(onCreate).toHaveBeenCalled())
    await userEvent.click(screen.getByRole('button', { name: /remove nature/i }))

    expect(onChange).toHaveBeenCalledWith([])
  })

  it('closes the dropdown but keeps the text when blurred during a pending create that then fails', async () => {
    let resolve!: (v: Item | null) => void
    const onCreate = vi.fn().mockReturnValue(new Promise<Item | null>((r) => { resolve = r }))
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={[nature]} onCreate={onCreate} placeholder="Add…" />)

    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'wip')
    expect(screen.getByText(/Create/)).toBeInTheDocument()
    await userEvent.keyboard('{Enter}')
    await userEvent.tab()

    await waitFor(() => expect(screen.queryByText(/Create/)).not.toBeInTheDocument())
    expect(input).toHaveValue('wip')

    resolve(null)
    await waitFor(() => expect(onCreate).toHaveBeenCalledTimes(1))
    expect(screen.queryByText(/Create/)).not.toBeInTheDocument()
    expect(input).toHaveValue('wip')
  })
})

describe('TokenInput — focus list cap', () => {
  const many = (n: number): Item[] =>
    Array.from({ length: n }, (_, i) => ({ id: `i-${i + 1}`, name: `Item${String(i + 1).padStart(2, '0')}` }))
  const HINT = /type to search/i

  it('shows the first 10 and a hint row when 25 items exist', async () => {
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={many(25)} openOnFocus placeholder="Add…" />)

    await userEvent.click(screen.getByPlaceholderText('Add…'))

    const rows = await screen.findAllByRole('listitem')
    expect(rows).toHaveLength(11)
    expect(rows.slice(0, 10).map((r) => r.textContent)).toEqual(many(10).map((i) => i.name))
    expect(rows[10]).toHaveTextContent(HINT)
  })

  it('shows no hint when 10 or fewer items exist', async () => {
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={many(10)} openOnFocus placeholder="Add…" />)

    await userEvent.click(screen.getByPlaceholderText('Add…'))

    expect(await screen.findAllByRole('listitem')).toHaveLength(10)
    expect(screen.queryByText(HINT)).not.toBeInTheDocument()
  })

  it('selected items do not consume the cap', async () => {
    const all = many(25)
    render(<TokenInput items={all.slice(0, 3)} onChange={vi.fn()} suggestions={all} openOnFocus />)

    await userEvent.click(screen.getByRole('textbox'))

    const rows = await screen.findAllByRole('listitem')
    expect(rows.slice(0, 10).map((r) => r.textContent)).toEqual(all.slice(3, 13).map((i) => i.name))
    expect(rows).toHaveLength(11)
  })

  it('typing finds an item beyond the first 10', async () => {
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={many(25)} openOnFocus placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), 'Item22')

    expect(await screen.findByText('Item22')).toBeInTheDocument()
    expect(screen.queryByText(HINT)).not.toBeInTheDocument()
  })

  it('typed matches are not capped', async () => {
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={many(25)} openOnFocus placeholder="Add…" />)

    await userEvent.type(screen.getByPlaceholderText('Add…'), 'Item')

    expect(await screen.findAllByRole('listitem')).toHaveLength(25)
  })

  it('arrow keys skip the hint row and Enter on it adds nothing', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} suggestions={many(11)} openOnFocus placeholder="Add…" />)
    await userEvent.click(screen.getByPlaceholderText('Add…'))
    await screen.findByText(HINT)

    // 10 rows: 10 presses reach the last item, the 11th wraps to the first.
    await userEvent.keyboard('{ArrowDown>10/}')
    await userEvent.keyboard('{ArrowDown}')
    await userEvent.keyboard('{ArrowUp}')
    await userEvent.keyboard('{Enter}')

    expect(onChange).toHaveBeenCalledTimes(1)
    expect(onChange).toHaveBeenCalledWith([expect.objectContaining({ name: 'Item10' })])
  })

  it('clicking the hint row adds nothing', async () => {
    const onChange = vi.fn()
    render(<TokenInput items={[]} onChange={onChange} suggestions={many(25)} openOnFocus placeholder="Add…" />)
    await userEvent.click(screen.getByPlaceholderText('Add…'))

    await userEvent.click(await screen.findByText(HINT))

    expect(onChange).not.toHaveBeenCalled()
  })

  it('clicking the hint row keeps the input focused and the list open', async () => {
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={many(25)} openOnFocus placeholder="Add…" />)
    const input = screen.getByPlaceholderText('Add…')
    await userEvent.click(input)

    await userEvent.click(await screen.findByText(HINT))
    await new Promise((r) => setTimeout(r, 200))

    expect(input).toHaveFocus()
    expect(screen.getByText(HINT)).toBeInTheDocument()
  })
})

describe('TokenInput — blur timer', () => {
  it('refocusing before the blur timer fires keeps the text and the open list', async () => {
    render(<TokenInput items={[]} onChange={vi.fn()} suggestions={[nature, travel]} placeholder="Add…" />)
    const input = screen.getByPlaceholderText('Add…')
    await userEvent.type(input, 'nat')
    expect(await screen.findByText('Nature')).toBeInTheDocument()

    await userEvent.tab()
    await userEvent.click(input)
    await new Promise((r) => setTimeout(r, 250))

    expect(input).toHaveValue('nat')
    expect(screen.getByText('Nature')).toBeInTheDocument()
  })
})
