package utils

import (
	"os"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
)

func TestPathExists(t *testing.T) {
	tmp, err := os.MkdirTemp("", "temp")
	if err != nil {
		t.Fatalf("An error occurred creating temp folder %s: %v", tmp, err)
	}

	defer os.RemoveAll(tmp)

	if !utils.PathExists(tmp) {
		t.Errorf("FAILED: got false, expected true for path %s", tmp)
	}
}

func TestPathDoesntExist(t *testing.T) {
	if utils.PathExists("temp") {
		t.Errorf("FALED, got true, expected false")
	}
}
