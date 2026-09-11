import ResourceCard from '@/components/ResourceCard'
import { type Character } from '../api/useCharacters'

interface CharactersGridProps {
  characters: Character[]
  onDeleteClick: (id: string) => void
  onEditClick: (character: Character) => void
}

export default function CharactersGrid({ characters, onDeleteClick, onEditClick }: CharactersGridProps) {
  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
      {characters.map((character) => (
        <ResourceCard
          key={character.id}
          imageUrl={character.avatar_url}
          label={character.name}
          onDeleteClick={() => onDeleteClick(character.id)}
          onClick={() => onEditClick(character)}
        />
      ))}
    </div>
  )
}
