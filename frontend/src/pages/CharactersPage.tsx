import { useState } from 'react'
import { Plus } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useCharacters } from '@/features/characters/api/useCharacters'
import CharactersGrid from '@/features/characters/components/CharactersGrid'

export default function CharactersPage() {
  const [search, setSearch] = useState('')
  const { data: characters = [], isLoading } = useCharacters()

  const filtered = search.trim()
    ? characters.filter((c) =>
        c.name.toLowerCase().includes(search.toLowerCase())
      )
    : characters

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Input
          placeholder="Search characters..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
        <Button variant="outline" className="ml-auto">
          <Plus />
          New
        </Button>
      </div>

      {isLoading ? null : <CharactersGrid characters={filtered} />}
    </div>
  )
}
