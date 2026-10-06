import { useState } from 'react'

export function useCellErrors() {
  const [cellErrors, setCellErrors] = useState<Record<string, string | null>>({})

  function setCellError(id: string, field: string, error: boolean) {
    setCellErrors(prev => ({ ...prev, [`${id}-${field}`]: error ? field : null }))
  }

  function hasCellError(id: string, field: string) {
    return cellErrors[`${id}-${field}`] != null
  }

  return { setCellError, hasCellError }
}
