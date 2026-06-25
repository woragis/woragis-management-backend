package service

import "strings"

const (
	maxUploadBytesImage = 10 << 20  // 10 MiB
	maxUploadBytesVideo = 100 << 20 // 100 MiB
)

func IsGalleryMime(mimeType string) bool {
	m := strings.ToLower(strings.TrimSpace(mimeType))
	return strings.HasPrefix(m, "image/") || strings.HasPrefix(m, "video/")
}

func maxBytesForMime(mimeType string) int {
	if strings.HasPrefix(strings.ToLower(mimeType), "video/") {
		return maxUploadBytesVideo
	}
	return maxUploadBytesImage
}

func extFromMime(mimeType string) string {
	switch strings.ToLower(mimeType) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	case "video/quicktime":
		return ".mov"
	case "application/pdf":
		return ".pdf"
	default:
		return ".bin"
	}
}
