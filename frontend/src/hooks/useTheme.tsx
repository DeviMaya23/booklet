import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'

const THEMES = ['warm', 'lumen', 'sunless'] as const

export type Theme = (typeof THEMES)[number]

const STORAGE_KEY = 'booklet-theme'

interface ThemeContextValue {
  theme: Theme
  setTheme: (theme: Theme) => void
}

const ThemeContext = createContext<ThemeContextValue | null>(null)

function readStoredTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY)
  return THEMES.find((theme) => theme === stored) ?? 'warm'
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>(readStoredTheme)

  useEffect(() => {
    document.documentElement.dataset.theme = theme
    localStorage.setItem(STORAGE_KEY, theme)
  }, [theme])

  return <ThemeContext.Provider value={{ theme, setTheme }}>{children}</ThemeContext.Provider>
}

// eslint-disable-next-line react-refresh/only-export-components -- provider + hook are co-located per design decision
export function useTheme(): ThemeContextValue {
  const context = useContext(ThemeContext)
  if (!context) throw new Error('useTheme must be used within a ThemeProvider')
  return context
}
