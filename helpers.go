package nim

import (
	"os"
)

const (
	// DefaultAddress is used if no other is specified.
	defaultServerAddress = ":3000"
)

// detectAddress
func detectAddress(addr ...string) string {
	if len(addr) > 0 {
		return addr[0]
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return defaultServerAddress
}
