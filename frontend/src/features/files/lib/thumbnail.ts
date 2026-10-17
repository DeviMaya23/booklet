const MAX_SIZE = 600;

const ALPHA_TYPES = new Set(["image/png", "image/gif", "image/webp"]);

export interface ThumbnailResult {
  blob: Blob;
  width: number;
  height: number;
}

export async function generateThumbnail(file: File): Promise<ThumbnailResult> {
  const bitmap = await createImageBitmap(file);
  const { width: srcW, height: srcH } = bitmap;

  const scale = Math.min(1, MAX_SIZE / Math.max(srcW, srcH));
  const width = Math.round(srcW * scale);
  const height = Math.round(srcH * scale);

  const canvas = new OffscreenCanvas(width, height);
  const ctx = canvas.getContext("2d");
  if (!ctx) throw new Error("Failed to acquire 2D rendering context");
  ctx.drawImage(bitmap, 0, 0, width, height);
  bitmap.close();

  const usePng = ALPHA_TYPES.has(file.type);
  const blob = await canvas.convertToBlob({
    type: usePng ? "image/png" : "image/jpeg",
    quality: usePng ? undefined : 0.9,
  });

  return { blob, width: srcW, height: srcH };
}
