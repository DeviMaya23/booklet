package mime

import (
	"regexp"
	"strings"
)

var unsafeFilenameChars = regexp.MustCompile(`[/\\:*?"<>|]+`)

func SanitizeFilename(s string) string {
	s = unsafeFilenameChars.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	if s == "" {
		return "artpiece"
	}
	return s
}

func MimeTypeToExt(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	case "application/pdf":
		return ".pdf"
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	default:
		return ".bin"
	}
}

func IsImage(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

// ThumbnailExt returns the file extension for a thumbnail generated from the
// given source mime type. Alpha-capable formats produce PNG; all others JPEG.
func ThumbnailExt(mimeType string) string {
	switch mimeType {
	case "image/png", "image/gif", "image/webp":
		return ".png"
	default:
		return ".jpg"
	}
}

// ThumbnailContentType returns the Content-Type for a thumbnail generated from
// the given source mime type.
func ThumbnailContentType(mimeType string) string {
	switch mimeType {
	case "image/png", "image/gif", "image/webp":
		return "image/png"
	default:
		return "image/jpeg"
	}
}
