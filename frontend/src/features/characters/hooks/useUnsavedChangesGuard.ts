import { useRef, useState } from 'react'
import { useBlocker } from 'react-router-dom'

interface UseUnsavedChangesGuardOptions {
  isDirty: boolean
  /** Called when the user leaves in place (no router navigation involved). */
  onLeave: () => void
}

/**
 * Guards an edit screen against losing unsaved changes. Covers both exits:
 * an in-page "back" (requestLeave) and in-app navigation (the router blocker).
 */
export function useUnsavedChangesGuard({ isDirty, onLeave }: UseUnsavedChangesGuardOptions) {
  const [leavePromptOpen, setLeavePromptOpen] = useState(false)
  const leavingRef = useRef(false)

  const blocker = useBlocker(() => !leavingRef.current && isDirty)
  const promptOpen = leavePromptOpen || blocker.state === 'blocked'

  function requestLeave() {
    if (isDirty) setLeavePromptOpen(true)
    else onLeave()
  }

  function discard() {
    setLeavePromptOpen(false)
    if (blocker.state === 'blocked') blocker.proceed()
    else onLeave()
  }

  function keepEditing() {
    setLeavePromptOpen(false)
    if (blocker.state === 'blocked') blocker.reset()
  }

  /** Call right before a navigation the user chose and the guard must not stop (after delete). */
  function allowLeaving() {
    leavingRef.current = true
  }

  return { promptOpen, requestLeave, discard, keepEditing, allowLeaving }
}
