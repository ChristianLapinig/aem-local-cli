package utils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ChristianLapinig/aem-local-cli/internal/utils"
)

func TestErrorAndCleanup(t *testing.T) {
	tmp := t.TempDir()
	testPath := filepath.Join(tmp, "test")
	if err := os.Mkdir(testPath, 0755); err != nil {
		t.Fatalf("Error creating folder: %v", err)
	}

	testErr := errors.New("test error")
	if err := utils.ErrorAndCleanup(testPath, testErr); err != nil && err.Error() != "test error" {
		t.Errorf("FAILED: got %s, expected 'test error'", err.Error())
	}
	if utils.PathExists(testPath) {
		t.Errorf("FAILED: expected %s to be deleted.", testPath)
	}
}
