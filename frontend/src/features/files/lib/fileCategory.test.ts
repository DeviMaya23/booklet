import { describe, it, expect } from 'vitest'
import { getFileCategory, getFileCategoryName } from './fileCategory'
import { Layers, Film, FileImage, File, FileArchive, FileText, Music, Box } from 'lucide-react'

describe('getFileCategory', () => {
  it('maps a known proprietary extension to the Layered category', () => {
    const result = getFileCategory('clip')
    expect(result.icon).toBe(Layers)
    expect(result.tinted).toBe(true)
    expect(result.cssVar).toBe('--file-cat-layered')
    expect(result.label).toBe('CLIP')
  })

  it('maps a video extension to Video category', () => {
    const result = getFileCategory('mp4')
    expect(result.icon).toBe(Film)
    expect(result.tinted).toBe(true)
    expect(result.cssVar).toBe('--file-cat-video')
    expect(result.label).toBe('MP4')
  })

  it('maps an image extension to Image category', () => {
    const result = getFileCategory('png')
    expect(result.icon).toBe(FileImage)
    expect(result.tinted).toBe(false)
    expect(result.cssVar).toBeNull()
    expect(result.label).toBe('PNG')
  })

  it('maps zip extension to Archive category', () => {
    const result = getFileCategory('zip')
    expect(result.icon).toBe(FileArchive)
    expect(result.label).toBe('ZIP')
  })

  it('maps rar extension to Archive category', () => {
    const result = getFileCategory('rar')
    expect(result.icon).toBe(FileArchive)
    expect(result.label).toBe('RAR')
  })

  it('maps docx extension to Document category', () => {
    const result = getFileCategory('docx')
    expect(result.icon).toBe(FileText)
    expect(result.label).toBe('DOCX')
  })

  it('maps vrm extension to 3D category', () => {
    const result = getFileCategory('vrm')
    expect(result.icon).toBe(Box)
    expect(result.tinted).toBe(true)
  })

  it('maps mp3 extension to Audio category', () => {
    const result = getFileCategory('mp3')
    expect(result.icon).toBe(Music)
    expect(result.tinted).toBe(true)
  })

  it('uppercases the extension for the label', () => {
    const result = getFileCategory('PSD')
    expect(result.icon).toBe(Layers)
    expect(result.label).toBe('PSD')
  })

  it('returns fallback with no label for null extension', () => {
    const result = getFileCategory(null)
    expect(result.icon).toBe(File)
    expect(result.tinted).toBe(false)
    expect(result.cssVar).toBeNull()
    expect(result.label).toBeNull()
  })

  it('returns fallback with label for unrecognised extension', () => {
    const result = getFileCategory('xyz')
    expect(result.icon).toBe(File)
    expect(result.label).toBe('XYZ')
  })
})

describe('getFileCategoryName', () => {
  it('returns the category name for a known extension', () => {
    expect(getFileCategoryName('mp4')).toBe('video')
    expect(getFileCategoryName('png')).toBe('image')
  })

  it('returns "other" for an unknown extension instead of "fallback"', () => {
    expect(getFileCategoryName('pptx')).toBe('other')
    expect(getFileCategoryName('xyz')).toBe('other')
  })

  it('returns "other" for null extension', () => {
    expect(getFileCategoryName(null)).toBe('other')
  })
})
