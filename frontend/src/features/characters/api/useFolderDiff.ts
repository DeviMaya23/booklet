import { useEffect, useState } from 'react'
import type { PublicFolder } from './usePublicFolders'

interface UseFolderDiffParams {
  originalFolders: { id: string; name: string }[]
  publicFolders: PublicFolder[] | undefined
}

export function useFolderDiff({ originalFolders, publicFolders }: UseFolderDiffParams) {
  const [addedFolders, setAddedFolders] = useState<PublicFolder[]>([])
  const [removedIds, setRemovedIds] = useState<Set<string>>(new Set())
  const [renamedFolders, setRenamedFolders] = useState<Map<string, string>>(new Map())
  const [invalidFolders, setInvalidFolders] = useState<{ id: string; name: string }[]>([])

  // Build a live token map for enriching selected folder rows
  const liveMap = new Map<string, PublicFolder>(
    (publicFolders ?? []).map((f) => [f.id, f]),
  )

  // Compute current selected folders: (original + added) - removed, enriched with live tokens
  const selectedFolders: PublicFolder[] = [
    ...originalFolders
      .filter((f) => !removedIds.has(f.id))
      .map((f) => ({
        id: f.id,
        name: renamedFolders.get(f.id) ?? f.name,
        token: liveMap.get(f.id)?.token ?? '',
      })),
    ...addedFolders,
  ]

  const selectedIds = new Set(selectedFolders.map((f) => f.id))

  useEffect(() => {
    if (!publicFolders) return
    const liveLookup = new Map<string, string>(publicFolders.map((f) => [f.id, f.name]))
    const nameOverrides = new Map<string, string>()
    const newRemovedIds = new Set<string>()
    const newInvalid: { id: string; name: string }[] = []
    for (const f of originalFolders) {
      const liveName = liveLookup.get(f.id)
      if (liveName === undefined) {
        newRemovedIds.add(f.id)
        newInvalid.push({ id: f.id, name: f.name })
      } else if (liveName !== f.name) {
        nameOverrides.set(f.id, liveName)
      }
    }
    if (nameOverrides.size > 0) {
      setRenamedFolders(nameOverrides)
    }
    if (newRemovedIds.size > 0) {
      setRemovedIds((prev) => new Set([...prev, ...newRemovedIds]))
      setInvalidFolders((prev) => {
        const existingIds = new Set(prev.map((f) => f.id))
        return [...prev, ...newInvalid.filter((f) => !existingIds.has(f.id))]
      })
    }
  // intentionally omits originalFolders: it's a new array reference every render;
  // this effect must only re-run when the server data changes, not on every render
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [publicFolders])

  function addFolder(folder: PublicFolder) {
    if (!selectedIds.has(folder.id)) {
      setAddedFolders((prev) => [...prev, folder])
    }
  }

  function removeFolder(id: string) {
    const isOriginal = originalFolders.some((f) => f.id === id)
    if (isOriginal) {
      setRemovedIds((prev) => new Set([...prev, id]))
    } else {
      setAddedFolders((prev) => prev.filter((f) => f.id !== id))
    }
  }

  function dismissInvalid(id: string) {
    setInvalidFolders((prev) => prev.filter((f) => f.id !== id))
  }

  function reset() {
    setAddedFolders([])
    setRemovedIds(new Set())
    setRenamedFolders(new Map())
    setInvalidFolders([])
  }

  return {
    selectedFolders,
    selectedIds,
    invalidFolders,
    addFolder,
    removeFolder,
    dismissInvalid,
    reset,
  }
}
