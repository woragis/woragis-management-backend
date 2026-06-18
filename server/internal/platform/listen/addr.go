package listen

import (
	"os"
	"strings"
)

// Addr returns the HTTP listen address.
// Railway and similar platforms set PORT; HTTP_ADDR overrides when set explicitly.
func Addr() string {
	if addr := strings.TrimSpace(os.Getenv("HTTP_ADDR")); addr != "" {
		return addr
	}
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.HasPrefix(port, ":") {
			return port
		}
		return ":" + port
	}
	return ":8080"
}
