import { useState } from 'react'

export function useDeleteAccountConfirmation(email: string | undefined) {
  const [confirming, setConfirming] = useState(false)
  const [typedEmail, setTypedEmail] = useState('')

  const emailMatches = !!email && typedEmail.trim().toLowerCase() === email.trim().toLowerCase()

  function start() {
    setConfirming(true)
  }

  function cancel() {
    setConfirming(false)
    setTypedEmail('')
  }

  function updateTypedEmail(value: string) {
    setTypedEmail(value)
  }

  return { confirming, typedEmail, emailMatches, start, cancel, updateTypedEmail }
}
