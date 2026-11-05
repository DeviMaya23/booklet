import { useParams } from 'react-router-dom'
import CharacterDetailView from '@/features/characters/components/CharacterDetailView'

export default function CharacterDetailPage() {
  const { id } = useParams<{ id: string }>()
  if (!id) return null
  return <CharacterDetailView key={id} characterId={id} />
}
