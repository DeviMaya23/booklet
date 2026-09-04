import { useState } from 'react'
import { Plus } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useImages } from '@/features/images/api/useImages'
import ImagesGrid from '@/features/images/components/ImagesGrid'

export default function ImagesPage() {
  const [search, setSearch] = useState('')
  const { data: images = [], isLoading } = useImages()

  const filtered = search.trim()
    ? images.filter(
        (img) =>
          img.title !== null &&
          img.title.toLowerCase().includes(search.toLowerCase())
      )
    : images

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Input
          placeholder="Search images..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
        <Button variant="outline" className="ml-auto">
          <Plus />
          New
        </Button>
      </div>

      {isLoading ? null : <ImagesGrid images={filtered} />}
    </div>
  )
}
