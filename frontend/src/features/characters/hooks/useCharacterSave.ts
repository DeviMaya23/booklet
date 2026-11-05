import { useState } from 'react'
import { useKindeAuth } from '@kinde-oss/kinde-auth-react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { apiFetch } from '@/lib/api'
import { CHARACTERS_QUERY_KEY, type Character } from '../api/useCharacters'
import { useUpdateCharacter } from '../api/useUpdateCharacter'
import { useInitAvatarUpload } from '../api/useInitAvatarUpload'
import { useCompleteAvatarUpload } from '../api/useCompleteAvatarUpload'
import { uploadCharacterAvatar } from '../lib/uploadCharacterAvatar'
import type { AvatarPicker } from './useAvatarPicker'

export interface CharacterFormValues {
  name: string
  notes: string
  verseId: string | null
  folderIds: string[]
}

interface UseCharacterSaveOptions {
  character: Character
  avatar: Pick<AvatarPicker, 'localFile' | 'avatarCleared' | 'isChanged' | 'reset'>
  onSaved: () => void
}

/**
 * Saves an edited character: the avatar change first (upload, or delete), then
 * the character update. Failures toast and leave the caller's edits untouched.
 */
export function useCharacterSave({ character, avatar, onSaved }: UseCharacterSaveOptions) {
  const { getToken } = useKindeAuth()
  const queryClient = useQueryClient()
  const updateMutation = useUpdateCharacter()
  const initUploadMutation = useInitAvatarUpload()
  const completeUploadMutation = useCompleteAvatarUpload()
  const [isSaving, setIsSaving] = useState(false)

  async function save({ name, notes, verseId, folderIds }: CharacterFormValues) {
    if (!name.trim() || isSaving) return
    setIsSaving(true)
    try {
      if (avatar.localFile) {
        try {
          await uploadCharacterAvatar(
            character.id,
            avatar.localFile,
            initUploadMutation.mutateAsync,
            completeUploadMutation.mutateAsync,
          )
        } catch {
          toast.error('Failed to upload avatar')
          return
        }
      } else if (avatar.avatarCleared) {
        const res = await apiFetch(`/characters/${character.id}/avatar`, getToken, { method: 'DELETE' })
        if (!res.ok) {
          toast.error('Failed to remove avatar')
          return
        }
      }
      if (avatar.isChanged) {
        // The avatar is saved from here on; show it even if the update below fails.
        avatar.reset()
        queryClient.invalidateQueries({ queryKey: CHARACTERS_QUERY_KEY })
      }

      await updateMutation.mutateAsync({
        id: character.id,
        name: name.trim(),
        notes: notes.trim() || null,
        folder_ids: folderIds,
        verse_id: verseId,
      })

      toast.success('Character updated')
      onSaved()
    } catch {
      toast.error('Failed to update character')
    } finally {
      setIsSaving(false)
    }
  }

  return { isSaving, save }
}
