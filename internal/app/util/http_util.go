package util

// FilterContentType extract mime type from header: Content-Type
func FilterContentType(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}
