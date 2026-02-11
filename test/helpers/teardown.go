package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
)

func Teardown(t testing.TB) {
	home := utils.GetHomePath()
	markerPath := filepath.Join(home, ".aemlocal_path")
	if err := os.RemoveAll(markerPath); err != nil {
		t.Fatalf("Error cleaning up marker file: %v", err)
	}
}
