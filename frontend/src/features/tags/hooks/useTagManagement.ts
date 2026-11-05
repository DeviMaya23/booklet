import { useState } from 'react'
import { toast } from 'sonner'
import { TagConflictError, type Tag } from '../api/useTags'
import { useUpdateTag } from '../api/useUpdateTag'
import { useDeleteTag } from '../api/useDeleteTag'
import { useMergeTag, type MergeTarget } from '../api/useMergeTag'
import { findTagByName } from '../lib/findTagByName'

// Orchestrates rename / merge / delete for the Tags page: which row is being
// edited, which confirmation is open, and which mutation a submit routes to.
export function useTagManagement(tags: Tag[]) {
  const [editingId, setEditingId] = useState<string | null>(null)
  const [renameError, setRenameError] = useState<string | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Tag | null>(null)
  const [mergeTarget, setMergeTarget] = useState<MergeTarget | null>(null)

  const updateTag = useUpdateTag()
  const deleteTag = useDeleteTag()
  const mergeTag = useMergeTag()

  function startEdit(tag: Tag) {
    setRenameError(null)
    setEditingId(tag.id)
  }

  function cancelEdit() {
    setRenameError(null)
    setEditingId(null)
  }

  async function submitRename(tag: Tag, name: string) {
    const collision = findTagByName(tags, name, tag.id)
    if (collision) {
      setMergeTarget({ source: tag, target: collision })
      return
    }
    setRenameError(null)
    try {
      await updateTag.mutateAsync({ id: tag.id, name })
      setEditingId(null)
    } catch (err) {
      setRenameError(
        err instanceof TagConflictError ? 'A tag with that name already exists' : 'Failed to rename tag',
      )
    }
  }

  function requestDelete(tag: Tag) {
    setDeleteTarget(tag)
  }

  function cancelDelete() {
    setDeleteTarget(null)
  }

  async function confirmDelete() {
    if (!deleteTarget) return
    try {
      await deleteTag.mutateAsync(deleteTarget.id)
      setDeleteTarget(null)
    } catch {
      toast.error('Failed to delete tag')
    }
  }

  // Cancelling leaves the row in edit mode so the name can be changed.
  function cancelMerge() {
    setMergeTarget(null)
  }

  async function confirmMerge() {
    if (!mergeTarget) return
    try {
      await mergeTag.mutateAsync({ sourceId: mergeTarget.source.id, targetId: mergeTarget.target.id })
      setMergeTarget(null)
      setEditingId(null)
    } catch (err) {
      setMergeTarget(null)
      toast.error(err instanceof Error ? err.message : 'Failed to merge tags')
    }
  }

  return {
    editingId,
    renameError,
    isRenaming: updateTag.isPending,
    deleteTarget,
    isDeleting: deleteTag.isPending,
    mergeTarget,
    isMerging: mergeTag.isPending,
    startEdit,
    cancelEdit,
    submitRename,
    requestDelete,
    cancelDelete,
    confirmDelete,
    cancelMerge,
    confirmMerge,
  }
}
