import {
  Film,
  Music,
  Layers,
  Box,
  FileArchive,
  FileText,
  FileImage,
  File,
  type LucideIcon,
} from 'lucide-react'

export type FileCategory = 'video' | 'audio' | 'layered' | '3d' | 'archive' | 'document' | 'image' | 'fallback'

export interface FileCategoryResult {
  icon: LucideIcon
  label: string | null
  tinted: boolean
  cssVar: string | null
}

const CATEGORY_META: Record<FileCategory, { icon: LucideIcon; tinted: boolean; cssVar: string | null }> = {
  video:    { icon: Film,        tinted: true,  cssVar: '--file-cat-video' },
  audio:    { icon: Music,       tinted: true,  cssVar: '--file-cat-audio' },
  layered:  { icon: Layers,      tinted: true,  cssVar: '--file-cat-layered' },
  '3d':     { icon: Box,         tinted: true,  cssVar: '--file-cat-3d' },
  archive:  { icon: FileArchive, tinted: false, cssVar: null },
  document: { icon: FileText,    tinted: false, cssVar: null },
  image:    { icon: FileImage,   tinted: false, cssVar: null },
  fallback: { icon: File,        tinted: false, cssVar: null },
}

// Extension → category mapping. Extension-first is the only signal; mime type is not used.
// Files uploaded before the extension field was added show the fallback icon (extension is null).
const EXTENSION_CATEGORY: Record<string, FileCategory> = {
  // Video
  mp4: 'video', mov: 'video', webm: 'video', mkv: 'video', m4v: 'video', avi: 'video',
  wmv: 'video', flv: 'video', ogv: 'video', ts: 'video',

  // Audio
  mp3: 'audio', wav: 'audio', flac: 'audio', m4a: 'audio', ogg: 'audio', aac: 'audio',
  wma: 'audio', opus: 'audio', aiff: 'audio', alac: 'audio',

  // Layered / project
  psd: 'layered', clip: 'layered', kra: 'layered', sai: 'layered', sai2: 'layered',
  xcf: 'layered', procreate: 'layered', ai: 'layered', afphoto: 'layered',
  psb: 'layered', csp: 'layered', afdesign: 'layered', pdn: 'layered',
  ase: 'layered', aseprite: 'layered',

  // 3D / rigged model
  blend: '3d', fbx: '3d', obj: '3d', glb: '3d', gltf: '3d', vrm: '3d',
  moc3: '3d', cmo3: '3d', vox: '3d', dae: '3d', stl: '3d',

  // Archive
  zip: 'archive', rar: 'archive', '7z': 'archive', tar: 'archive', gz: 'archive',
  tgz: 'archive', bz2: 'archive', xz: 'archive', zst: 'archive',

  // Document
  pdf: 'document', txt: 'document', md: 'document', doc: 'document', docx: 'document',
  rtf: 'document', odt: 'document', pages: 'document',

  // Image
  jpg: 'image', jpeg: 'image', png: 'image', gif: 'image', svg: 'image', webp: 'image',
  tiff: 'image', tif: 'image', bmp: 'image', heic: 'image', heif: 'image', avif: 'image',
  ico: 'image',
}

export function getFileCategory(extension: string | null): FileCategoryResult {
  const ext = extension?.toLowerCase() ?? null
  const category = (ext ? EXTENSION_CATEGORY[ext] : undefined) ?? 'fallback'
  const meta = CATEGORY_META[category]
  const label = ext ? ext.toUpperCase() : null
  return { ...meta, label }
}
