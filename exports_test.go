package nim

import "testing"

func TestNimRun(t *testing.T) {
	// just test that Run doesn't bomb
	go New().Run(":3000")
}

func TestNimDefault(t *testing.T) {
	go Default().Run(":3000")
}
