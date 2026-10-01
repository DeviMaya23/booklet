package mime

import "strings"

func MimeTypeToExt(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}

func IsImage(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}
