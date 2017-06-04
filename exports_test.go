package nim

import (
	"os"
	"testing"
)

func TestNimRun(t *testing.T) {
	// just test that Run doesn't bomb
	go Run(New(), ":3000")
}

func TestNimDefault(t *testing.T) {
	go Run(Default(), ":3000")
}

func TestDetectAddress(t *testing.T) {
	if detectAddress() != defaultServerAddress {
		t.Error("Expected the defaultServerAddress")
	}

	if detectAddress(":6060") != ":6060" {
		t.Error("Expected the provided address")
	}

	os.Setenv("PORT", "8080")
	if detectAddress() != ":8080" {
		t.Error("Expected the PORT env var with a prefixed colon")
	}
}
