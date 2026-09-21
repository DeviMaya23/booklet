import { useState } from 'react'
import { Plus } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useImages, type Image } from '@/features/images/api/useImages'
import ImagesGrid from '@/features/images/components/ImagesGrid'
import ImageFormModal from '@/features/images/components/ImageFormModal'

export default function ImagesPage() {
  const [search, setSearch] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editingImage, setEditingImage] = useState<Image | undefined>(undefined)
  const { data: images = [], isLoading } = useImages()

  const filtered = search.trim()
    ? images.filter(
        (img) =>
          img.title !== null &&
          img.title.toLowerCase().includes(search.toLowerCase())
      )
    : images

  function openCreate() {
    setEditingImage(undefined)
    setModalOpen(true)
  }

  function openEdit(image: Image) {
    setEditingImage(image)
    setModalOpen(true)
  }

  function handleModalOpenChange(open: boolean) {
    setModalOpen(open)
    if (!open) setEditingImage(undefined)
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Input
          placeholder="Search images..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
        <Button variant="outline" className="ml-auto" onClick={openCreate}>
          <Plus />
          New
        </Button>
      </div>

      {isLoading ? null : <ImagesGrid images={filtered} onEditClick={openEdit} />}

      <ImageFormModal
        open={modalOpen}
        onOpenChange={handleModalOpenChange}
        image={editingImage}
      />
    </div>
  )
}
